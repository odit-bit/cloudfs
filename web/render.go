package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/web/component"
)

func (v *App) serviceErr(w http.ResponseWriter, r *http.Request, err error) {
	v.logger.Error(fmt.Sprintf("View: %v\n", err), "path", r.URL.Path)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (v *App) writeObject(w http.ResponseWriter, r *http.Request, obj *blob.ObjectInfo) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", obj.Filename))

	defer obj.Data.Close()
	if _, err := io.Copy(w, obj.Data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
