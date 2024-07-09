package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/odit-bit/cloudfs/internal/ui"
)

func (app *api) serveUploadPage(serviceEndpoint string) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.getUserID(r)
		if err != nil {
			if err != nil {
				//redirect into loginHTML
				log.Printf("%s:%s \n", uploadAPIEndpoint, err)
				http.Redirect(w, r, loginHTML, http.StatusFound)
				return
			}
		}
		if err := ui.RenderUploadPage(w, serviceEndpoint); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
	return handler
}

func (app *api) handleUpload(w http.ResponseWriter, r *http.Request) {
	userID, err := app.getUserID(r)
	if err != nil {
		if err != nil {
			//redirect into loginHTML
			log.Printf("%s:%s \n", uploadAPIEndpoint, err)
			http.Redirect(w, r, loginHTML, http.StatusFound)
			return
		}
	}

	// parse file from request
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(r.RemoteAddr, err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileType, err := getFileType(header)
	if err != nil {
		log.Println("ulalal: ", err)
	}
	if fileType == "" {
		fileType = "application/octetstream"
	}

	// save blob into storage
	res, err := app.blobStorage.PutObject(r.Context(), userID, header.Filename, file, header.Size, minio.PutObjectOptions{
		ContentType: fileType,
	})

	if err != nil {
		log.Printf("%v: %v \n", uploadAPIEndpoint, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, res.ETag)
}

// =====================

var bufPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 512)
		return &buf
	},
}

func getFileType(filePart *multipart.FileHeader) (string, error) {
	file, err := filePart.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Seek back to the beginning of the file after detecting the content type
	_, _ = file.Seek(0, 0)

	// Use http.DetectContentType to guess the MIME type without reading the entire file
	buffer := *bufPool.Get().(*[]byte) // Adjust the read size as needed
	_, _ = file.Read(buffer)

	// Use DetectContentType
	fileType := http.DetectContentType(buffer)

	clear(buffer)
	bufPool.Put(&buffer)
	return fileType, nil
}
