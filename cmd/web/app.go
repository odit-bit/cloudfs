package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/minio/minio-go/v7"
	"github.com/odit-bit/cloudfs/internal/ui"
	"github.com/odit-bit/cloudfs/internal/user/pguser"
)

const (
	listHTML     = "/list"
	loginHTML    = "/login"
	registerHTML = "/register"
	uploadHTML   = "/upload"
)

const (
	listAPIEndpoint     = "/api/list"
	loginAPIEndpoint    = "/api/login"
	registerAPIEndpoint = "/api/register"
	uploadAPIEndpoint   = "/api/upload"
)

const (
	defaultMaxAge = 60 * 5
)

type api struct {
	blobStorage *minio.Client
	userDB      *pguser.DB //user.Database
	cs          *sessions.CookieStore
}

func (app *api) serveIndex(w http.ResponseWriter, r *http.Request) {
	idx := ui.NewIndexPage()
	idx.AddMenu("list", listHTML)
	idx.AddMenu("upload", uploadHTML)
	// idx.AddMenu("register", registerHTML)

	if err := idx.Render(w); err != nil {
		log.Println("serve index: ", err)
		return
	}
}

// get user id from session-cookie request
func (app *api) getUserID(r *http.Request) (string, error) {
	sess, err := app.cs.Get(r, "cloudfs_session")
	if err != nil {
		log.Println("get user id: ", err)
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
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   defaultMaxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	sess.Values["userID"] = id
	if err = sess.Save(r, w); err != nil {
		log.Println("saving session: ", err)
		return err
	}
	return nil
}
