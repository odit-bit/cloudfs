package api

import (
	"io"
	"net/http"
)

func (v *App) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	info, ok := getTokenCtx(ctx)
	if !ok {
		v.logger.Error("apiHandler: wrong userID context")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	filename := r.URL.Query().Get("filename")
	obj, err := v.objects.Object(ctx, info.UserID(), filename)
	if err != nil {
		v.serviceErr(w, r, "download", err)
		return
	}
	io.Copy(w, obj.Data)
}

// presigned token , non-authorized

func (v *App) DownloadWithToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")
	obj, err := v.objects.DownloadToken(ctx, token)
	if err != nil {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	io.Copy(w, obj.Data)
}
