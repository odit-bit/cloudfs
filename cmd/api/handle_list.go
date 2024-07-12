package main

import (
	"log"
	"net/http"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/ui"
)

const ()

func (app *api) serveListPage(serviceEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ui.RenderListPage(w, serviceEndpoint); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

// render html page with injected list data
func (app *api) HandleList(w http.ResponseWriter, r *http.Request) {
	// validate session
	userID, err := app.getUserIDFromCookie(r)
	if err != nil {
		//redirect into loginHTML
		log.Printf("%s:%s \n", listAPIEndpoint, err)
		http.Redirect(w, r, loginHTML, http.StatusFound)
		return
	}

	ctx := r.Context()
	c := app.blobStorage.List(ctx, userID, 100, "")
	var list []*blob.ObjectInfo

	for c.Next() {
		var info blob.ObjectInfo
		c.Scan(&info)
		list = append(list, &info)
	}

	if err := ui.RenderListResult(w, list, "/api/download"); err != nil {
		log.Println("render list result html: ", err)
		http.Error(w, err.Error(), 500)
	}

}
