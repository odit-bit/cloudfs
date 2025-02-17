package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) Route() *chi.Mux {
	mux := chi.NewMux()

	mux.Use(middleware.Logger)

	mux.Post("/api/v1/auth", a.Auth)
	mux.Post("/api/v1/register", a.Register)
	mux.Get("/api/v1/{share-token}", a.DownloadWithToken)

	mux.Group(func(r chi.Router) {
		r.Use(a.NeedToken)
		r.Get("/api/v1", a.Download)
		r.Post("/api/v1", a.Upload)
		r.Get("/api/v1/share", a.CreateShareToken)
	})

	return mux
}
