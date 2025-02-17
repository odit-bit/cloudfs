package blob

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_upload_chunk(t *testing.T) {

	storage, err := NewWithMemory()
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("this is data")
	size := len(data)
	_ = size
	res, err := storage.Upload(context.Background(), UploadParam{
		Bucket:      "bucket",
		Filename:    "filename",
		ContentType: "content",
		Size:        int64(size),
		Body:        bytes.NewBuffer(data),
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(data), int(res.Size))

}
