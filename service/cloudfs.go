package service

import (
	"context"
	"io"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/user"
)

type BlobStore interface {
	Put(ctx context.Context, bucket, filename string, reader io.Reader, size int64, contentType string) (*blob.ObjectInfo, error)
	Get(ctx context.Context, bucket, filename string) (*blob.ObjectInfo, error)
	Delete(ctx context.Context, bucket, filename string) error
	ObjectIterator(ctx context.Context, bucket string, limit int, lastFilename string) *blob.Iterator
	// GenerateShareURL(ctx context.Context, bucket, filename string) (*url.URL, error)
}

type TokenStore interface {
	// Validate(ctx context.Context, tokenString string) (userID, filename string, ok bool)
	// Generate(ctx context.Context, bucket, filename string, dur time.Duration) (string, error)
	Query(ctx context.Context, txn func(txn TokenTxn) error) error
}
type TokenTxn interface {
	Put(ctx context.Context, token *blob.ShareToken) error
	Get(ctx context.Context, tokenString string) (*blob.ShareToken, error)
	Delete(ctx context.Context, key string) error
	Commit() error
	Cancel() error
}

type AccountStore interface {
	Find(ctx context.Context, username string) (*user.Account, error)
	Insert(ctx context.Context, acc *user.Account) error
}

type Cloudfs struct {
	blobService    BlobStore
	tokenService   TokenStore
	accountService AccountStore
}

func NewCloudfs(tokenStore TokenStore, blobStore BlobStore, accountStore AccountStore) (*Cloudfs, error) {
	return &Cloudfs{
		blobService:    blobStore,
		tokenService:   tokenStore,
		accountService: accountStore,
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
