package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dustin/go-humanize"
	"github.com/odit-bit/cloudfs/component"
	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/user"
)

type App struct {
	logger *slog.Logger
	// svc      *service.Cloudfs
	accounts *user.Users
	objects  *blob.Blobs
	session  *scs.SessionManager
}

func New(users *user.Users, obj *blob.Blobs, sess *scs.SessionManager, logger *slog.Logger) *App {
	ah := App{
		logger:   logger,
		accounts: users,
		objects:  obj,
		session:  sess,
	}
	return &ah
}

func (app *App) Run(ctx context.Context, addr string, middlewares ...func(http.Handler) http.Handler) error {
	handler := app.router(middlewares...)
	srv := http.Server{
		Addr:        addr,
		Handler:     handler,
		IdleTimeout: 10 * time.Second,
	}
	// setup listener
	lAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", lAddr)
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		v, ok := <-sig
		if !ok {
			panic("unexpected closed signal ")
		}
		close(sig)
		app.logger.Info("shutdown server", "type", v)
		srv.Close()
	}(ctx)

	app.logger.Info(fmt.Sprintf("listen on %s", l.Addr().String()))
	if err := srv.Serve(l); err != nil {
		if err != http.ErrServerClosed {
			app.logger.Error(err.Error())
		} else {
			err = nil
		}
	}
	if lErr := l.Close(); lErr != nil {
		err = errors.Join(err, lErr)
	}

	wg.Wait()
	app.logger.Info("gracefully shutdown")
	return err
}

// redirect with http status 303
func (v *App) LoginService(redirecURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.PostFormValue("username")
		pass := r.PostFormValue("password")
		acc, err := v.accounts.BasicAuth(r.Context(), user, pass)
		if err != nil {
			v.logger.Error(fmt.Sprintf("authHandler: %v", err))
			http.Error(w, "username not found", http.StatusNotFound)
			return
		}

		v.session.Put(r.Context(), "userID", acc.ID.String())
		http.Redirect(w, r, redirecURL, http.StatusSeeOther)
	}
}

// return midlleware-like handler that place in front of handler that needed for auth
func (v *App) auth(loginRedirectURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := v.session.GetString(r.Context(), "userID")
			if userID == "" {
				v.logger.Error("session return invalid userID")
				http.Redirect(w, r, loginRedirectURL, http.StatusSeeOther)
				return
			}

			ctx := contextSetUser(r.Context(), (*UserID)(&userID))
			rr := r.WithContext(ctx)
			next.ServeHTTP(w, rr)
		})
	}
}

func (v *App) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	UserID, ok := getUserIDFromCtx(ctx)
	if !ok {
		v.logger.Error("apiHandler: wrong userID context")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	filename := r.URL.Query().Get("filename")
	obj, err := v.objects.Object(ctx, UserID, filename)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	v.writeObject(w, r, obj)
}

func (v *App) DownloadShare(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")
	obj, err := v.objects.DownloadToken(ctx, token)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	v.writeObject(w, r, obj)
}

func (v *App) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := getUserIDFromCtx(ctx)
	if !ok {
		v.logger.Error("apiHandler: wrong userID context")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	lastFilename := r.URL.Query().Get("last")

	objects, err := v.objects.ListObject(r.Context(), userID, 100, lastFilename)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}

	attr := component.ListAttribute{
		Objects:      objects,
		DownloadAPI:  _DownloadService,
		DeleteAPI:    _DeleteService,
		ListView:     _ListView,
		ShareFileAPI: _ShareFileService,
	}

	listDataView := component.ListData(&attr)
	listDataView.Render(r.Context(), w)

}

// redirect with http status 303
func (v *App) RegisterService(redirectURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username := r.FormValue("username")
		if username == "" {
			http.Error(w, "username cannot be nil", http.StatusBadRequest)
			return
		}
		pass := r.FormValue("password")
		if pass == "" {
			http.Error(w, "password cannot be nil", http.StatusBadRequest)
			return
		}

		//insert account
		if err := v.accounts.Register(r.Context(), username, pass); err != nil {
			if err == user.ErrAccountExist {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			v.logger.Error(err.Error())
			http.Error(w, "failed create account", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}

}

func (v *App) ShareFile(publicDownloadPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := getUserIDFromCtx(ctx)
		if !ok {
			v.logger.Error("apiHandler: wrong userID context")
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
		filename := r.URL.Query().Get("filename")
		tkn, err := v.objects.CreateShareToken(r.Context(), userID, filename, 0)
		if err != nil {
			v.serviceErr(w, r, err)
			return
		}

		// w.Header().Set("Path", obj.Token)
		q := url.Values{}
		q.Set("token", tkn.Key())
		shareURL := url.URL{
			// Scheme:   r.URL.Scheme,
			Host:     r.Host,
			Path:     publicDownloadPath,
			RawQuery: q.Encode(),
		}

		comp := component.ShareFileResponse(shareURL.String(), humanize.Time(tkn.ValidUntil()))
		comp.Render(ctx, w)
	}
}

func (v *App) Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Trigger", "newObject")
	defer r.Body.Close()
	//get userID
	ctx := r.Context()
	userID, _ := getUserIDFromCtx(ctx)
	if userID == "" {
		v.logger.Error("apiHandler: userID is empty")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	fd, err := handleMultipart(r, "file")
	if err != nil {
		v.logger.Error("upload handler failed parse multipart", "err", err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer fd.Close()
	result, err := v.objects.Upload(r.Context(), blob.UploadParam{
		Bucket:      userID,
		Filename:    fd.Filename,
		Size:        fd.Size,
		ContentType: fd.ContentType,
		Body:        fd.Body,
	})
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	_ = result
	w.WriteHeader(http.StatusOK)
}

func (v *App) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	UserID, ok := getUserIDFromCtx(ctx)
	if !ok {
		v.logger.Error("apiHandler: wrong userID context")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	filename := r.URL.Query().Get("filename")
	err := v.objects.Delete(ctx, UserID, filename)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	w.Header().Set("HX-Trigger", "deleteObject")
	w.WriteHeader(http.StatusOK)
}
