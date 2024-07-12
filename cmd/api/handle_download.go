package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/odit-bit/cloudfs/internal/blob"
)

func (app *api) handleDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := app.getUserIDFromCookie(r)
	if err != nil {
		log.Println("download handler: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filename := r.URL.Query().Get("filename")
	var objInfo blob.ObjectInfo
	err = app.blobStorage.Get(ctx, userID, filename, &objInfo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer objInfo.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", objInfo.ObjName))
	w.Header().Set("Content-Type", objInfo.ContentType)

	if _, err := io.Copy(w, objInfo.Reader()); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
