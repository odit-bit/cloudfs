package web

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dustin/go-humanize"
	"github.com/odit-bit/cloudfs/server/apipb"
	"github.com/odit-bit/cloudfs/web/component"
	"github.com/odit-bit/cloudfs/internal/storage"
	"github.com/odit-bit/cloudfs/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type App struct {
	logger *slog.Logger
	// svc      *service.Cloudfs
	// accounts *user.Users
	// objects *storage.Blobs
	session *scs.SessionManager

	backend apipb.StorageServiceClient
}

func New(backendAddr string, sess *scs.SessionManager, logger *slog.Logger) (*App, error) {
	creds := insecure.NewCredentials()
	conn, err := grpc.NewClient(backendAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	backend := apipb.NewStorageServiceClient(conn)
	ah := App{
		logger:  logger,
		session: sess,
		backend: backend,
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
		acc, err := v.backend.BasicAuth(r.Context(), &apipb.BasicAuthRequest{Username: user, Password: pass})
		if err != nil {
			v.logger.Error(fmt.Sprintf("authHandler: %v", err))
			http.Error(w, "username not found", http.StatusNotFound)
			return
		}

		v.session.Put(r.Context(), sessUserToken, acc.Token)
		v.logger.Info(fmt.Sprintf("loginService put token: %v", acc.Token))
		http.Redirect(w, r, redirecURL, http.StatusSeeOther)
	}
}

// return midlleware-like handler that place in front of handler that needed for auth
func (v *App) auth(loginRedirectURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := v.session.GetString(r.Context(), sessUserToken)
			if token == "" {
				v.logger.Error("session return invalid userToken")
				http.Redirect(w, r, loginRedirectURL, http.StatusSeeOther)
				return
			}

			ctx := setUserTokenCtx(r.Context(), (*userToken)(&token))
			rr := r.WithContext(ctx)
			next.ServeHTTP(w, rr)
		})
	}
}

func (v *App) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userToken, ok := getUserTokenFromCtx(ctx)
	if !ok {
		v.logger.Error("apiHandler: wrong userID context")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	filename := r.URL.Query().Get("filename")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))

	objStream, err := v.backend.DownloadObject(ctx, &apipb.DownloadRequest{
		Token:    userToken,
		Filename: filename,
	})
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	defer objStream.CloseSend()

	for {
		res, err := objStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			v.serviceErr(w, r, err)
			return
		}
		w.Write(res.Chunk)
	}
	// v.writeObject(w, r)
}

func (v *App) DownloadShare(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("shareToken")

	stream, err := v.backend.DownloadSharedObject(ctx, &apipb.DownloadSharedRequest{SharedToken: token})
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}

	md, err := stream.Header()
	if err != nil {
		v.logger.Error(fmt.Sprintf("failed receive header: %v", err))
		v.serviceErr(w, r, err)
		return
	}
	if md == nil {
		_, xerr := stream.Recv()
		err = errors.Join(err, xerr)
		v.serviceErr(w, r, err)
		return
	}
	var filename, contentType string
	if xfilename := md.Get("filename"); len(xfilename) == 0 {
		v.serviceErr(w, r, fmt.Errorf("missing header filename from server"))
		return
	} else {
		filename = xfilename[0]
	}
	if xct := md.Get("filename"); len(xct) == 0 {
		v.serviceErr(w, r, fmt.Errorf("missing header filename from server"))
		return
	} else {
		contentType = xct[0]
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))
	w.Header().Set("Content-Type", contentType)

	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			v.serviceErr(w, r, errors.Join(err, stream.CloseSend()))
			return
		}
		w.Write(res.Chunk)
	}

	// v.writeObject(w, r)
}

func (v *App) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userToken, ok := getUserTokenFromCtx(ctx)
	if !ok {
		v.logger.Error("apiHandler: wrong userID context")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	lastFilename := r.URL.Query().Get("last")

	stream, err := v.backend.ListObject(r.Context(), &apipb.ListObjectRequest{
		UserToken:    userToken,
		Limit:        1000,
		LastFilename: lastFilename,
	})
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	defer stream.CloseSend()

	// CANDIDATE TO HTMX SSE
	var objects []storage.ObjectInfo

	for {
		obj, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			v.serviceErr(w, r, err)
			return
		}
		objects = append(objects, storage.ObjectInfo{
			Bucket:       obj.UserID,
			Filename:     obj.Filename,
			ContentType:  obj.ContentType,
			Sum:          obj.Sum,
			Size:         obj.Size,
			LastModified: obj.LastModified.AsTime(),
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
		if _, err := v.backend.Register(r.Context(), &apipb.RegisterRequest{
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

func (v *App) ShareFile(publicDownloadPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		shareToken, ok := getUserTokenFromCtx(ctx)
		if !ok {
			v.logger.Error("apiHandler: wrong userID context")
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
		filename := r.URL.Query().Get("filename")
		tkn, err := v.backend.ShareObject(r.Context(), &apipb.ShareObjectRequest{
			Token:    shareToken,
			Filename: filename,
		})
		if err != nil {
			v.serviceErr(w, r, err)
			return
		}

		// w.Header().Set("Path", obj.Token)
		q := url.Values{}
		q.Set("shareToken", tkn.ShareToken)
		shareURL := url.URL{
			// Scheme:   r.URL.Scheme,
			Host:     r.Host,
			Path:     publicDownloadPath,
			RawQuery: q.Encode(),
		}

		comp := component.ShareFileResponse(shareURL.String(), humanize.Time(tkn.ValidUntil.AsTime().UTC()))
		comp.Render(ctx, w)
	}
}

func (v *App) Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Trigger", "newObject")
	defer r.Body.Close()
	//get userID
	ctx := r.Context()
	userToken, _ := getUserTokenFromCtx(ctx)
	if userToken == "" {
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

	md := metadata.New(map[string]string{})
	md.Set("filename", fd.Filename)
	md.Set("authorization", userToken)
	md.Set("content-type", fd.ContentType)
	md.Set("content-length", strconv.Itoa(int(fd.Size)))

	sCtx := metadata.NewOutgoingContext(ctx, md)
	stream, err := v.backend.UploadObject(sCtx)
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	defer stream.CloseSend()

	chunk := make([]byte, 1024*1024*4)
	for {
		n, err := fd.Body.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := stream.Send(&apipb.UploadRequest{
			Token:       userToken,
			Filename:    fd.Filename,
			TotalSize:   fd.Size,
			ContentType: fd.ContentType,
			Chunk:       chunk[:n]},
		); err != nil {
			v.logger.Error("failed write chunk to backend")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (v *App) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userToken, ok := getUserTokenFromCtx(ctx)
	if !ok {
		v.logger.Error("apiHandler: wrong userID context")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	filename := r.URL.Query().Get("filename")
	_, err := v.backend.DeleteObject(ctx, &apipb.DeleteRequest{
		UserToken: userToken,
		Filename:  filename,
	})
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}

	w.Header().Set("HX-Trigger", "deleteObject")
	w.WriteHeader(http.StatusOK)
}
