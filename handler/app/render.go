package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/odit-bit/cloudfs/component"
	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/service"
)

func (v *App) serviceErr(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, service.ErrUpload) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, service.ErrBucketNotExisted) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	v.logger.Error(fmt.Sprintf("View: %v\n", err), "path", r.URL.Path)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (v *App) writeObject(w http.ResponseWriter, r *http.Request, obj *blob.ObjectInfo) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", obj.Filename))
	reader, err := obj.Reader.Get(r.Context())
	if err != nil {
		v.serviceErr(w, r, err)
		return
	}
	defer reader.Close()

	if _, err := io.Copy(w, reader); err != nil {
		v.serviceErr(w, r, err)
		return
	}
}

func loginPage(loginAPI, registerPageURL string) http.HandlerFunc {
	h := component.Login(loginAPI, registerPageURL)
	return func(w http.ResponseWriter, r *http.Request) {
		h.Render(r.Context(), w)
	}
}

func registerPage(registerAPI string) http.HandlerFunc {
	h := component.Register(registerAPI)
	return func(w http.ResponseWriter, r *http.Request) {
		h.Render(r.Context(), w)
	}
}

func indexPage(uploadAPI, listViewURL string) http.HandlerFunc {
	idxPage := component.Index(uploadAPI, listViewURL)
	return func(w http.ResponseWriter, r *http.Request) {
		idxPage.Render(r.Context(), w)

	}
}
