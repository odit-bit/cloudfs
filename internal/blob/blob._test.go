package blob

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_iterator(t *testing.T) {

	tc := []struct {
		input []*ObjectInfo
	}{
		{input: []*ObjectInfo{
			{
				UserID:   "1",
				Filename: "obj1",
			},
			{
				UserID:   "2",
				Filename: "obj2",
			},
		}},
		{input: []*ObjectInfo{}},
	}

	// //
	// infos := make([]*ObjectInfo, 2)
	// infos[0] = &ObjectInfo{
	// 	UserID:   "1",
	// 	Filename: "obj1",
	// }
	// infos[1] = &ObjectInfo{
	// 	UserID:   "2",
	// 	Filename: "obj2",
	// }

	for _, test := range tc {
		objC := make(chan *ObjectInfo, 10)

		for _, v := range test.input {
			objC <- v
		}
		close(objC)

		it := Iterator{
			UserID: "",
			C:      objC,
			obj:    &ObjectInfo{},
			// err:    nil,
		}

		for _, expected := range test.input {
			if !it.Next() {
				t.Fatal("should return info")
			}

			obj := it.Value()
			assert.Equal(t, expected.UserID, obj.UserID)
		}

	}
}

func Test_readBlob(t *testing.T) {

	data := []byte("this is data")

	buf := bytes.NewBuffer(data)
	rc := io.NopCloser(buf)
	obj := ObjectInfo{
		// UserID:       "123",
		// Filename:     "file-123",
		// ContentType:  "secret",
		// Sum:          "",
		// Size:         int64(len(data)),
		// LastModified: time.Time{},
		Reader: ReaderFunc(func(ctx context.Context) (io.ReadCloser, error) {
			return rc, nil
		}),
		// isReaded: false,
	}

	dst := bytes.Buffer{}
	src, err := obj.Reader.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	n, err := io.Copy(&dst, src)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, n, obj.Size)
	assert.Equal(t, data, dst.Bytes())
}
