package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/ui"
	"github.com/odit-bit/cloudfs/internal/user/pguser"
)

type api struct {
	blobStorage *blob.Storage //*minio.Client
	userDB      *pguser.DB    //user.Database
	cs          *sessions.CookieStore
}

func NewApiHandler(conf *config) *api {
	cmdCtx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	//setup user database
	udb, err := pguser.NewDB(
		cmdCtx,
		conf.Storage.User.URI,
	)

	if err != nil {
		log.Fatal(err)
	}

	//setup storage
	storage, err := blob.NewStorage(
		conf.Storage.Blob.Endpoint,
		conf.Storage.Blob.AccessKey,
		conf.Storage.Blob.SecretKey,
	)
	if err != nil {
		log.Fatal(err)
	}

	// setup cookie store
	//TODO: move it to own package as component
	cs := sessions.NewCookieStore([]byte(conf.Session.Token))
	cs.Options = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   60 * 5,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	api := api{
		blobStorage: storage,
		userDB:      udb,
		cs:          cs,
	}

	return &api
}

func (app *api) serveIndex(w http.ResponseWriter, _ *http.Request) {
	idx := ui.NewIndexPage()
	idx.AddMenu("list", listHTML)
	idx.AddMenu("upload", uploadHTML)
	// idx.AddMenu("register", registerHTML)

	if err := idx.Render(w); err != nil {
		log.Println("serve index: ", err)
		return
	}
}
