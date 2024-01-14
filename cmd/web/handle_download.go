package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
)

func (app *api) handleDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := app.getUserID(r)
	if err != nil {
		log.Println("download handler: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filename := r.URL.Query().Get("filename")
	// ri, err := app.bs.Retreive(ctx, userID, filename)
	ri, err := app.blobStorage.GetObject(ctx, userID, filename, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	stat, err := ri.Stat()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", stat.Key))
	w.Header().Set("Content-Type", stat.ContentType)

	if _, err := io.Copy(w, ri); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
