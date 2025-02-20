package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_upload_chunk(t *testing.T) {

	storage, err := NewWithMemory()
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("this is data")
	sum := sum256(data)
	size := len(data)
	_ = size
	param := UploadParam{
		Bucket:      "bucket",
		Filename:    "filename",
		ContentType: "content",
		Size:        int64(size),
		// Body:        bytes.NewBuffer(data),
	}
	cw := storage.NewChunkWriter(
		context.Background(),
		param.Bucket,
		param.Filename,
		param.ContentType,
		param.Size,
	)

	start := 0
	end := 1
	for end <= len(data) {
		if _, err := cw.Write(data[start:end]); err != nil {
			t.Fatal(err)
		}
		start = end
		end++
	}

	res, err := cw.Result()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, sum, res.Sum)
	assert.Equal(t, len(data), int(res.Size))

}

func sum256(p []byte) string {
	hash := sha256.New()
	hash.Write(p)
	return hex.EncodeToString(hash.Sum(nil))
}
