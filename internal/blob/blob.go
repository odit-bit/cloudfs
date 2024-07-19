package blob

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type ShareFunc func(ctx context.Context, expiration time.Duration) (*url.URL, error)

func (f ShareFunc) GetURL(ctx context.Context, expiration time.Duration) (*url.URL, error) {
	return f(ctx, expiration)
}

//---------------

type ReaderFunc func(ctx context.Context) (io.ReadCloser, error)

func (rf ReaderFunc) Get(ctx context.Context) (io.ReadCloser, error) {
	return rf(ctx)
}

type Reader interface {
	Get(ctx context.Context) (io.ReadCloser, error)
}

//---------------

// represent object
type ObjectInfo struct {
	UserID       string
	Filename     string
	ContentType  string
	Sum          string
	Size         int64
	LastModified time.Time

	// //shareURL
	// ShareFn ShareFunc
	//get object content
	Reader Reader
}

func (o *ObjectInfo) Validate() error {
	if o.Reader == nil {
		return fmt.Errorf("object reader is nil")
	}
	return nil
}

/////////// cursor

type Cursor struct {
	UserID  string
	listC   <-chan minio.ObjectInfo
	objInfo *minio.ObjectInfo

	err error
}

func (li *Cursor) Next() bool {
	objInfo, ok := <-li.listC
	if !ok {
		return false
	}
	if objInfo.Err != nil {
		li.err = objInfo.Err
		return false
	}

	li.objInfo = &objInfo
	return true
}

func (li *Cursor) Scan(info *ObjectInfo) {
	if info == nil {
		panic("storage cursor cannot scan into nil info")
	}

	// fmt.Println("STORAGE CURSOR", li.objInfo)
	info.UserID = li.UserID
	info.Filename = li.objInfo.Key
	info.Sum = li.objInfo.ETag
	info.Size = li.objInfo.Size
	info.LastModified = li.objInfo.LastModified

	li.objInfo = nil
}

func (li *Cursor) Error() error {
	return li.err
}

//////////

type Iterator struct {
	UserID string
	C      <-chan *ObjectInfo
	obj    *ObjectInfo
	err    error
}

func (li *Iterator) Next() bool {
	info, ok := <-li.C
	if !ok {
		return false
	}

	li.obj = info
	return true
}

func (li *Iterator) Value() *ObjectInfo {
	return li.obj
}

func (li *Iterator) Error() error {
	return li.err
}
