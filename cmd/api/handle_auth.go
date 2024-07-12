package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/odit-bit/cloudfs/internal/ui"
	"github.com/odit-bit/cloudfs/internal/user"
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

////

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

	//TODO: make it transaction
	//insert account
	if err := app.userDB.Insert(r.Context(), acc); err != nil {
		log.Println(err)
		http.Error(w, "failed create account", http.StatusInternalServerError)
		return
	}

	// make bucket
	err := app.blobStorage.CreateBucket(r.Context(), userID)
	if err != nil {
		err = fmt.Errorf("%v: %v", registerAPIEndpoint, err)
		log.Printf("%v \n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sess, err := app.cs.Get(r, "cloudfs_session")
	if err != nil {
		log.Println("get session: ", err)
	}

	sess.Values["userID"] = userID
	if err = sess.Save(r, w); err != nil {
		log.Println("saving session: ", err)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// get user id from session-cookie request
func (app *api) getUserIDFromCookie(r *http.Request) (string, error) {
	sess, err := app.cs.Get(r, "cloudfs_session")
	if err != nil {
		log.Println("get session: ", err)
	}

	userID, ok := sess.Values["userID"].(string)
	if userID == "" || !ok {
		return "", fmt.Errorf("not login")
	}

	return userID, nil
}

// save session into w (cookie)
func (app *api) saveSession(id string, w http.ResponseWriter, r *http.Request) error {
	sess, err := app.cs.Get(r, "cloudfs_session")
	if err != nil {
		log.Println("save session: ", err)
	}

	sess.Values["userID"] = id
	if err = sess.Save(r, w); err != nil {
		log.Println("saving session: ", err)
		return err
	}
	return nil
}
