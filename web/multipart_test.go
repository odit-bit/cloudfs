package web

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handle_multpart(t *testing.T) {
	urlString := "http://localhost:8181/api/upload"
	payload := strings.NewReader("-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"file\"; filename=\"Cintaku-1.mp3\"\r\nContent-Type: audio/mpeg\r\n\r\nthis is content\r\n-----011000010111000001101001--\r\n")
	req, _ := http.NewRequest("PUT", urlString, payload)
	req.Header.Add("Content-Type", "multipart/form-data; boundary=---011000010111000001101001")
	req.Header.Add("Content-length", "165")

	fi, err := handleMultipart(req, "file")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()
	assert.Equal(t, int64(165), fi.Size)
	assert.Equal(t, "audio/mpeg", fi.ContentType)
	assert.Equal(t, "Cintaku-1.mp3", fi.Filename)
}
