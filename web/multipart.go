package web

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
)

type fileInfo struct {
	Size        int64
	ContentType string
	Filename    string
	Body        io.ReadCloser
}

func (f *fileInfo) Close() error {
	return f.Body.Close()
}

func handleMultipart(req *http.Request, formName string) (*fileInfo, error) {
	fileSize := req.Header.Get("X-File-Size")
	if fileSize == "" {
		return nil, fmt.Errorf("invalid 'X-File-Size' value %v", fileSize)
	}
	size, err := strconv.Atoi(req.Header.Get("X-File-Size"))
	if err != nil {
		return nil, err
	}
	if size <= 0 {
		return nil, fmt.Errorf("file size cannot be nil, this is a bug")
	}
	// filename := req.Header.Get("X-File-Name")
	// ct := req.Header.Get("X-File-Type")

	mt, param, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if mt != "multipart/form-data" {
		return nil, fmt.Errorf("not mulitpart request")
	}

	reader := multipart.NewReader(req.Body, param["boundary"])
	part, err := reader.NextPart()
	if err != nil {
		return nil, err
	}

	ct := part.Header.Get("content-type")
	cd, param, err := mime.ParseMediaType(part.Header.Get("content-disposition"))
	if err != nil {
		return nil, err
	}
	if cd != "form-data" {
		return nil, fmt.Errorf("not form-data request")
	}

	name := param["name"]
	if name != formName {
		return nil, fmt.Errorf("wrong form name, got '%s' expect 'file' ", name)
	}
	filename := param["filename"]

	return &fileInfo{
		Size:        int64(size),
		ContentType: ct,
		Filename:    filename,
		Body:        part,
	}, nil
}
