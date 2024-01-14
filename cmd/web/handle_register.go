package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/odit-bit/cloudfs/internal/ui"
	"github.com/odit-bit/cloudfs/internal/user"
)

func (app *api) serveRegisterPage(serviceEndpoint string) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := ui.RenderRegisterPage(w, serviceEndpoint); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})

	return handler
}

func (app *api) handleRegister(w http.ResponseWriter, r *http.Request) {

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

	acc := user.CreateAccount(username, pass)
	userID := acc.ID.String()
	err := app.blobStorage.MakeBucket(r.Context(), userID, minio.MakeBucketOptions{
		Region:        "",
		ObjectLocking: false,
	})
	if err != nil {
		err = fmt.Errorf("%v: %v", registerAPIEndpoint, err)
		log.Printf("%v \n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := app.userDB.Insert(r.Context(), acc); err != nil {
		log.Println(err)
		http.Error(w, "failed create account", http.StatusInternalServerError)
		return
	}

	// log.Printf("%v: success create bucket %v \n", registerAPIEndpoint, userID)
	if err := app.saveSession(userID, w, r); err != nil {
		printErr(w, err, registerAPIEndpoint)
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
