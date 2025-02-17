package api

import (
	"fmt"
	"net/http"

	"github.com/odit-bit/cloudfs/internal/blob"
)

type uploadHeader struct {
	filename    string
	size        int64
	contentType string
}

func validateUploadHeader(r *http.Request) (*uploadHeader, error) {
	size := r.ContentLength
	if size <= 0 {
		return nil, fmt.Errorf("invalid size object")
	}
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		return nil, fmt.Errorf("illegal filename")
	}
	ct := r.Header.Get("Content-Type")
	return &uploadHeader{filename: filename, size: size, contentType: ct}, nil

}

func (v *App) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	info, ok := getTokenCtx(ctx)
	if !ok {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	uh, err := validateUploadHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if res, err := v.objects.Upload(ctx, blob.UploadParam{
		Bucket:      info.UserID(),
		Filename:    uh.filename,
		Size:        uh.size,
		ContentType: uh.contentType,
		Body:        r.Body,
	}); err != nil {
		v.serviceErr(w, r, "upload", err)
		return
	} else {
		fmt.Fprintf(w, "result: %v \n", res)
	}
}
