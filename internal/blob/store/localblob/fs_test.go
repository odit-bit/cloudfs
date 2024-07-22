package localblob

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_(t *testing.T) {
	v, err := New("./zzz/temp")
	if err != nil {
		t.Fatal(err)
	}
	defer v.purge()

	type obj struct {
		bucket   string
		filename string
		data     []byte
	}

	input := obj{
		bucket:   "user-id-1",
		filename: "my-file",
		data:     []byte("content-file-1"),
	}

	//Put
	info, err := v.Put(context.TODO(), input.bucket, input.filename, bytes.NewReader(input.data), int64(len(input.data)), "")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, input.filename, info.Filename)

	//Get
	actual, err := v.Get(context.Background(), input.bucket, input.filename)
	if err != nil {
		t.Fatal(err)
	}
	rc, err := actual.Reader.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()

	data := bytes.Buffer{}
	io.Copy(&data, rc)
	assert.Equal(t, string(input.data), data.String())

	//delete
	if err := v.Delete(context.Background(), input.bucket, input.filename); err != nil {
		t.Fatal(err)
	}
}

func Test_iterator(t *testing.T) {
	v, err := New("./zzz/temp2")
	if err != nil {
		t.Fatal(err)
	}
	defer v.purge()

	type obj struct {
		bucket   string
		filename string
		data     []byte
	}

	bucket := "my-bucket"
	objects := []obj{
		{
			bucket:   bucket,
			filename: "my-file-1",
			data:     []byte("123"),
		},
		{
			bucket:   bucket,
			filename: "my-file-2",
			data:     []byte("123"),
		},
	}
	for _, obj := range objects {
		_, err := v.Put(context.Background(), obj.bucket, obj.filename, bytes.NewReader(obj.data), int64(len(obj.data)), "")
		if err != nil {
			t.Fatal(err)
		}
	}

	iter := v.ObjectIterator(context.Background(), bucket, 100, "")
	if iter.Error() != nil {
		t.Fatal(err)
	}

	count := 0
	for iter.Next() {
		info := iter.Value()
		assert.Equal(t, objects[count].filename, info.Filename)
		count++
	}
	assert.Equal(t, len(objects), 2)

	iter2 := v.ObjectIterator(context.Background(), bucket, 1, objects[0].filename)
	if iter2.Error() != nil {
		t.Fatal(err)
	}

	if ok := iter2.Next(); !ok {
		t.Fatal(ok)
	}
	info := iter2.Value()
	assert.Equal(t, objects[1].filename, info.Filename)

}
