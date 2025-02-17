package blob

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_afero(t *testing.T) {
	ctx := context.Background()
	v, _ := newObjectMemory() //repo.NewAferoBlob("")
	_ = v
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
	info, err := v.Put(ctx, input.bucket, input.filename, bytes.NewReader(input.data), int64(len(input.data)), "")
	if err != nil {
		t.Fatal(err)
	}

	//Get
	actual, err := v.Get(ctx, input.bucket, input.filename)
	if err != nil {
		t.Fatal(err)
	}
	defer actual.Data.Close()
	assert.Equal(t, input.filename, info.Filename)
	data, _ := io.ReadAll(actual.Data)
	assert.Equal(t, input.data, data)

	// list
	iter := v.ObjectIterator(ctx, input.bucket, 1000, "")
	list := []ObjectInfo{}
	for obj := range iter.C {
		list = append(list, obj)
	}
	if len(list) == 0 {
		t.Fatal("list length should 1")
	}
	assert.Equal(t, input.filename, list[0].Filename)

	//delete
	if err := v.Delete(context.Background(), input.bucket, input.filename); err != nil {
		t.Fatal(err)
	}
}

func Test_local_iterator(t *testing.T) {
	v, err := newObjectMemory()
	if err != nil {
		t.Fatal(err)
	}

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
