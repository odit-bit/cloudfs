package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/odit-bit/cloudfs/component"
	"github.com/odit-bit/cloudfs/service"
)

type App struct {
	logger  *slog.Logger
	svc     *service.Cloudfs
	session *scs.SessionManager
}

func New(svc *service.Cloudfs, sess *scs.SessionManager, logger *slog.Logger) *App {
	ah := App{
		logger:  logger,
		svc:     svc,
		session: sess,
	}
	return &ah
}

func (app *App) Run(ctx context.Context, addr string, middlewares ...func(http.Handler) http.Handler) error {
	handler := app.router(middlewares...)
	srv := http.Server{
		Addr:    addr,
		Handler: handler,
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
		ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx2)
	}(ctx)

	app.logger.Info(fmt.Sprintf("listen on %s", addr))
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		app.logger.Error(err.Error())
	} else {
		err = nil
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
		acc, err := v.svc.Auth(r.Context(), &service.AuthParam{
			Username: user,
			Password: pass,
		})
		if err != nil {
			v.logger.Error(fmt.Sprintf("authHandler: %v", err))
			http.Error(w, "username not found", http.StatusNotFound)
			return
		}

		v.session.Put(r.Context(), "userID", acc.ID)
		http.Redirect(w, r, redirecURL, http.StatusSeeOther)
	}
}

// return midlleware like handler that place in front of handler that needed for auth
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
	obj, err := v.svc.Object(ctx, UserID, filename)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", obj.Filename))

	rc, err := obj.Reader.Get(ctx)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	defer rc.Close()

	if _, err := io.Copy(w, rc); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
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
	err := v.svc.Delete(ctx, UserID, filename)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	w.Header().Set("HX-Trigger", "deleteObject")
	w.WriteHeader(http.StatusOK)
}

func (v *App) Upload(w http.ResponseWriter, r *http.Request) {
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
		v.logger.Error(err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer fd.Close()
	obj, err := v.svc.Upload(r.Context(), &service.UploadParam{
		UserID:      userID,
		Filename:    fd.Filename,
		Size:        -1,
		ContentType: fd.ContentType,
		DataReader:  fd.Body,
	})
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	_ = obj
	// objs := []*blob.ObjectInfo{obj}
	// listView(objs, _ListView, w, r)
	w.Header().Set("HX-Trigger", "newObject")
	w.WriteHeader(http.StatusOK)
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

	objects, err := v.svc.ListObject(r.Context(), userID, 100, lastFilename)
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
		if err := v.svc.Register(r.Context(), &service.RegisterParam{
			Username: username,
			Password: pass,
		}); err != nil {
			if err == service.ErrAccountExist {
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
		obj, err := v.svc.SharingFile(r.Context(), userID, filename)
		if err != nil {
			v.serviceErr(w, r, err)
			return
		}

		// w.Header().Set("Path", obj.Token)
		q := url.Values{}
		q.Set("token", obj.Token)
		shareURL := url.URL{
			// Scheme:   r.URL.Scheme,
			Host:     r.Host,
			Path:     publicDownloadPath,
			RawQuery: q.Encode(),
		}

		comp := component.ShareFileResponse(shareURL.String(), obj.ValidUntil)
		comp.Render(ctx, w)
	}
}

func (v *App) DownloadShare(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.URL.Query().Get("token")
	if err := v.svc.DownloadSharedFile(ctx, token, func(r io.Reader) {
		io.Copy(w, r)
	}); err != nil {
		v.serviceErr(w, r, err)
		return
	}
}
