package main

import (
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
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
	userID, err := app.getUserID(r)
	if err != nil {
		//redirect into loginHTML
		log.Printf("%s:%s \n", listAPIEndpoint, err)
		http.Redirect(w, r, loginHTML, http.StatusFound)
		return
	}

	c := app.blobStorage.ListObjects(r.Context(), userID, minio.ListObjectsOptions{})
	ctx := r.Context()
	var list []*ui.ListData
	for {
		select {
		case <-ctx.Done():
		case info, ok := <-c:
			if ok {
				if info.Err == nil {
					ld := ui.ListData{
						Name:      info.Key,
						Size:      info.Size,
						SharedURL: "",
					}
					list = append(list, &ld)
					continue
				}

			}
		}
		break
	}

	if err := ui.RenderListResult(w, list, "/api/download"); err != nil {
		log.Println("render list result html: ", err)
		http.Error(w, err.Error(), 500)
	}

}
