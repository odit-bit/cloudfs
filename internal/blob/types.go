package blob

import (
	"context"
	"io"
	"net/url"
	"time"
)

type ShareFunc func(ctx context.Context, expiration time.Duration) (*url.URL, error)

func (f ShareFunc) GetURL(ctx context.Context, expiration time.Duration) (*url.URL, error) {
	return f(ctx, expiration)
}

// represent object
type Info struct {
	Bucket       string
	Filename     string
	ContentType  string
	Sum          string
	Size         int64
	LastModified time.Time
	// Data io.ReadCloser
}

type Data io.ReadCloser

type Object struct {
	Info
	Data
}

//////////

type Iterator struct {
	UserID string
	C      <-chan *Info
	obj    *Info
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

func (li *Iterator) Value() *Info {
	return li.obj
}

func (li *Iterator) Err() error {
	return li.err
}

//////

type QuotaInfo struct {
	Bucket string
	Quota  int64
	Usage  int64
}

func (qi *QuotaInfo) CheckAvail(size int64) bool {
	after := qi.Usage + size
	return after <= qi.Quota
}
