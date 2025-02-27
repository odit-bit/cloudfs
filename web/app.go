package web

import (
	"context"
	"errors"
	"fmt"
	"io"
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
	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/blob/blobpb"
	"github.com/odit-bit/cloudfs/internal/user"
	"github.com/odit-bit/cloudfs/internal/user/userpb"
	"github.com/odit-bit/cloudfs/web/component"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	logger *logrus.Logger
	// svc      *service.Cloudfs
	// accounts *user.Users
	// objects *storage.Blobs
	session *scs.SessionManager

	backend *Backend
}

func New(backendAddr string, sess *scs.SessionManager, logger *logrus.Logger) (*App, error) {
	creds := insecure.NewCredentials()
	conn, err := grpc.NewClient(backendAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	backend := Backend{
		auth:    userpb.NewAuthServiceClient(conn),
		objects: blobpb.NewStorageServiceClient(conn),
	}
	ah := App{
		logger:  logger,
		session: sess,
		backend: &backend,
	}

	return &ah, nil
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
		srv.SetKeepAlivesEnabled(false)
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
		res, err := v.backend.BasicAuth(r.Context(), user, pass)
		if err != nil {
			v.logger.Error(fmt.Sprintf("authHandler: %v", err))
			http.Error(w, "username not found", http.StatusNotFound)
			return
		}

		acc := &account{UserID: res.UserID, Token: res.Token, SharedObjects: map[Filename]ShareToken{}}
		v.session.Put(r.Context(), Session_Account, acc)
		// v.logger.Info(fmt.Sprintf("loginService put token: %v", acc))
		http.Redirect(w, r, redirecURL, http.StatusSeeOther)
	}
}

// return midlleware-like handler that place in front of handler that needed for auth
func (a *App) auth(loginRedirectURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := a.session.Get(r.Context(), Session_Account)
			acc, ok := v.(*account)
			if !ok {
				a.logger.Debug("session return invalid account (unauthorized): %T", v)
				http.Redirect(w, r, loginRedirectURL, http.StatusSeeOther)
				return
			}

			ctx := setAccountCtx(r.Context(), acc)
			rr := r.WithContext(ctx)
			next.ServeHTTP(w, rr)
		})
	}
}

func (v *App) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	acc, ok := getAccountFromCtx(ctx)
	if !ok {
		v.logger.Error("apiHandler: wrong userID context")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	filename := r.URL.Query().Get("filename")
	res, err := v.backend.DownloadObject(ctx, acc.UserID, filename)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	defer res.Reader.Close()
	v.attachment(w, r, filename, res.Reader)
}

func (v *App) ShareFile(publicDownloadPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		acc, ok := getAccountFromCtx(ctx)
		if !ok {
			v.logger.Error("apiHandler: wrong userID context")
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		filename := r.URL.Query().Get("filename")
		shareToken, ok := acc.SharedObjects[Filename(filename)]
		// v.logger.Infof("SHARE_FILE token: %v, filename: %v \n", shareToken.Key, filename)
		if !ok || shareToken.IsExpired() {
			res, err := v.backend.ShareObject(r.Context(), acc.UserID, filename)
			if err != nil {
				v.serviceErr(w, r, err)
				return
			}
			shareToken.Key = res.ShareToken
			shareToken.ValidUntil = res.ValidUntil
			acc.SharedObjects[Filename(filename)] = shareToken
			v.session.Put(ctx, Session_Account, acc)
		}

		// w.Header().Set("Path", obj.Token)
		q := url.Values{}
		q.Set("shareToken", shareToken.Key)
		shareURL := url.URL{
			// Scheme:   r.URL.Scheme,
			Host:     r.Host,
			Path:     publicDownloadPath,
			RawQuery: q.Encode(),
		}

		comp := component.ShareFileResponse(shareToken.Key, shareURL.String(), humanize.Time(shareToken.ValidUntil.UTC()))
		comp.Render(ctx, w)
	}
}

func (v *App) DownloadShare(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shareToken := r.URL.Query().Get("shareToken")

	object, err := v.backend.DownloadWithToken(ctx, shareToken)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	defer object.Reader.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", object.Filename))
	w.Header().Set("Content-Type", object.ContentType)
	if _, err := io.Copy(w, object.Reader); err != nil {
		v.serviceErr(w, r, err)
		return
	}

	// v.writeObject(w, r)
}

func (v *App) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	acc, ok := getAccountFromCtx(ctx)
	if !ok {
		v.logger.Error("apiHandler: wrong userID context")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	defer v.session.Put(ctx, Session_Account, acc)

	lastFilename := r.URL.Query().Get("last")

	c, err := v.backend.Objects(r.Context(), acc.UserID, lastFilename, 1000)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}

	// CANDIDATE TO HTMX SSE
	var objects []blob.ObjectInfo
	for obj := range c {
		if obj.Err() != nil {
			v.serviceErr(w, r, obj.Err())
			break
		}
		objects = append(objects, blob.ObjectInfo{
			Bucket:       obj.Bucket,
			Filename:     obj.Filename,
			ContentType:  obj.ContentType,
			Sum:          obj.Sum,
			Size:         obj.Size,
			LastModified: obj.LastModified,
		})
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
		if _, err := v.backend.Register(r.Context(), RegisterParam{
			Username: username,
			Password: pass,
		}); err != nil {
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

func (v *App) Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Trigger", "newObject")
	defer r.Body.Close()

	//get userID
	ctx := r.Context()
	acc, _ := getAccountFromCtx(ctx)
	// defer v.session.Put(ctx, Session_Account, acc)

	fd, err := handleMultipart(r, "file")
	if err != nil {
		v.logger.Error("upload handler failed parse multipart", "err", err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer fd.Close()
	// v.logger.Infof("http header X-File-Size: %d", fd.Size)

	res, err := v.backend.UploadObject(r.Context(), acc.UserID, fd.Filename, fd.ContentType, fd.Size, fd.Body)
	if err != nil {
		v.serviceErr(w, r, fmt.Errorf("uploading: %v", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res.Sum))

}

func (v *App) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	acc, ok := getAccountFromCtx(ctx)
	if !ok {
		v.logger.Error("apiHandler: wrong userID context")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	defer v.session.Put(ctx, Session_Account, acc)

	filename := r.URL.Query().Get("filename")
	_, err := v.backend.Delete(ctx, acc.UserID, filename)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}

	w.Header().Set("HX-Trigger", "deleteObject")
	w.WriteHeader(http.StatusOK)
}
