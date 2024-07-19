package service

import (
	"context"
	"io"
	"sync"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/user"
)

type BlobStore interface {
	Put(ctx context.Context, bucket, filename string, reader io.Reader, size int64, contentType string) (*blob.ObjectInfo, error)
	Get(ctx context.Context, bucket, filename string) (*blob.ObjectInfo, error)
	Delete(ctx context.Context, bucket, filename string) error
	// GenerateShareURL(ctx context.Context, bucket, filename string) (*url.URL, error)
}

type BucketStore interface {
	CreateBucket(ctx context.Context, bucket string) (any, error)
	IsBucketExist(ctx context.Context, bucket string) (bool, error)
	ObjectIterator(ctx context.Context, bucket string, limit int, lastFilename string) *blob.Iterator
}

type AccountStore interface {
	Find(ctx context.Context, username string) (*user.Account, error)
	Insert(ctx context.Context, acc *user.Account) error
}

type Cloudfs struct {
	blobService    BlobStore
	bucketService  BucketStore
	accountService AccountStore

	mx     sync.Mutex
	biller map[string]float64
}

func NewCloudfs(bucketStore BucketStore, blobStore BlobStore, accountStore AccountStore) (*Cloudfs, error) {
	return &Cloudfs{
		blobService:    blobStore,
		bucketService:  bucketStore,
		accountService: accountStore,
		mx:             sync.Mutex{},
		biller:         map[string]float64{},
	}, nil
}

// func (app *Cloudfs) calculateBills(id string, n uint64) {
// 	price := 0.00000012 // per mb
// 	app.mx.Lock()
// 	defer app.mx.Unlock()
// 	bill := app.biller[id]
// 	bill += price * float64(n/humanize.MiByte)
// 	log.Printf("current bill %f \n", bill)
// 	app.biller[id] = bill
// }
