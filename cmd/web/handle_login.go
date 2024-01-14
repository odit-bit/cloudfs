package main

import (
	"log"
	"net/http"

	"github.com/odit-bit/cloudfs/internal/ui"
)

func (app *api) serveLoginPage(serviceEndpoint string) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := ui.RenderLoginPage(w, serviceEndpoint); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})

	return handler
}

func (app *api) handleLogin(w http.ResponseWriter, r *http.Request) {

	user := r.PostFormValue("username")
	acc, err := app.userDB.Find(r.Context(), user)
	if err != nil {
		log.Printf("%s: %s \n", loginAPIEndpoint, err)
		http.Error(w, "username not found", http.StatusNotFound)
		return
	}

	pass := r.PostFormValue("password")
	if ok := acc.CheckPassword(pass); !ok {
		log.Printf("%s: %s \n", loginAPIEndpoint, err)
		http.Error(w, "wrong user or password", http.StatusBadRequest)
		return
	}

	if err := app.saveSession(acc.ID.String(), w, r); err != nil {
		log.Printf("%s: %s \n", loginAPIEndpoint, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	//should it redirect from server or client side ??
	http.Redirect(w, r, listAPIEndpoint, http.StatusFound)

}
