package blob

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
)

func Test_storage(t *testing.T) {
	key := "admin"
	secret := "admin12345"

	bs, err := NewStorage(DefaultEndpoint, key, secret)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer([]byte("testinput"))
	userID := "12345"
	objName := "testInput1"
	contentType := "application/octet-stream"
	size := int64(buf.Len())
	payload := buf

	if err := bs.MakeBucket(context.Background(), userID); err != nil {
		t.Fatal(err)
	}

	if err := testUploadFile(bs, userID, objName, contentType, size, payload); err != nil {
		t.Fatal(err)
	}

	ri, err := testRetrieve(bs, userID, objName)
	if err != nil {
		t.Fatal(err)
	}
	if ri.ContentType != contentType {
		t.Fatalf("got %v, expect %v \n", ri.ContentType, contentType)
	}

	if err := testDeleteBLOB(bs, userID); err != nil {
		t.Fatal(err)
	}
}

func testUploadFile(bs *Storage, userID, objName, contentType string, size int64, reader io.Reader) error {

	sum, err := bs.Save(context.Background(), userID, objName, contentType, size, reader)
	if err != nil {
		return err
	}
	if sum == "" {
		return fmt.Errorf("save success but no sum returned")
	}

	return nil
}

func testRetrieve(bs *Storage, userID, objName string) (*RetreiveInfo, error) {
	ri, err := bs.Retreive(context.Background(), userID, objName)
	if err != nil {
		return nil, err
	}

	return ri, nil
}

func testDeleteBLOB(bs *Storage, userID string) error {
	return bs.purge(userID)
}
