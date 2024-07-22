package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	_health = "/api/health"

	_LoginService          = "/api/login"
	_RegisterService       = "/api/register"
	_UploadService         = "/api/upload"
	_DownloadService       = "/api/download"
	_DeleteService         = "/api/delete"
	_ShareFileService      = "/api/share"
	_PublicDownloadService = "/public/download"

	_LoginPage    = "/login"
	_RegisterPage = "/register"

	_ListView = "/list/view"
)

// implement http.Handler
func (v *App) router(middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middlewares...)
	mux.Use(v.session.LoadAndSave)

	// health
	mux.Get(_health, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index page"))
	})

	//authentication
	mux.Get(_RegisterPage, registerPage(_RegisterService))
	mux.Get(_LoginPage, loginPage(_LoginService, _RegisterPage))
	mux.Post(_LoginService, v.LoginService("/"))
	mux.Post(_RegisterService, v.RegisterService(_LoginPage))

	// public download link
	mux.Get(_PublicDownloadService, v.DownloadShare)

	// view group
	mux.Group(func(authGroup chi.Router) {
		authGroup.Use(v.auth(_LoginPage))

		authGroup.Get("/", indexPage(_UploadService, _ListView))
		authGroup.Get(_DownloadService, v.Download)
		authGroup.Get(_ShareFileService, v.ShareFile(_PublicDownloadService))

		authGroup.Post(_UploadService, v.Upload)
		authGroup.Delete(_DeleteService, v.Delete)
		authGroup.Get(_ListView, v.List)
	})

	return mux
}
