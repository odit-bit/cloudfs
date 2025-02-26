package blob

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidShareToken = errors.New("invalid share token")
	// ErrUnknown           = &opErr{err: errors.New("exception")}
)

type OpErr interface {
	error
	mustEmbedded()
}

type ObjectStorer interface {
	Put(ctx context.Context, bucket, filename string, reader io.ReadCloser, size int64, contentType string) (*ObjectInfo, error)
	Get(ctx context.Context, bucket, filename string) (*ObjectInfo, error)
	Delete(ctx context.Context, bucket, filename string) error
	ObjectIterator(ctx context.Context, bucket string, limit int, lastFilename string) *Iterator
}

type TokenStorer interface {
	Put(ctx context.Context, token *Token) OpErr
	Get(ctx context.Context, tokenKey string) (*Token, bool, error)
	GetByFilename(ctx context.Context, bucket string) (*Token, bool, error)
	Delete(ctx context.Context, tokenKey string) error
}

type Storage struct {
	tokens  TokenStorer
	objects ObjectStorer
}

func NewWithMemory() (*Storage, error) {
	objects, _ := newObjectMemory()
	tokens, _ := newObjectTokenMemory()
	return New(context.Background(), tokens, objects)
}

func New(ctx context.Context, tokens TokenStorer, objects ObjectStorer) (*Storage, error) {
	return &Storage{tokens: tokens, objects: objects}, nil
}

func (o *Storage) Download(ctx context.Context, bucket, filename string) (*ObjectInfo, error) {
	return o.objects.Get(ctx, bucket, filename)
}

type UploadParam struct {
	Bucket      string
	Filename    string
	Size        int64
	Body        io.ReadCloser
	ContentType string
}

func (o *Storage) Put(ctx context.Context, param *UploadParam) (*ObjectInfo, error) {
	obj, err := o.objects.Put(ctx, param.Bucket, param.Filename, param.Body, param.Size, param.ContentType)
	if err != nil {
		return nil, err
	}

	return obj, nil

}

func (o *Storage) CreateShareToken(ctx context.Context, bucket, filename string, expire time.Duration) (*Token, error) {
	// var tkn *Token
	if bucket == "" {
		return nil, fmt.Errorf("invalid bucket")
	}
	tkn, ok, _ := o.tokens.GetByFilename(ctx, filename)
	if !ok {
		if expire <= 60*time.Minute {
			expire = 60 * time.Minute
		}
		tkn := NewShareToken(bucket, filename, expire)
		err := o.tokens.Put(ctx, tkn)
		return tkn, err
	}
	return tkn, nil
}

func (o *Storage) DownloadToken(ctx context.Context, tknKey string) (*ObjectInfo, error) {
	tkn, ok, err := o.tokens.Get(ctx, tknKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidShareToken
	}
	if !tkn.IsNotExpire() {
		err2 := o.tokens.Delete(ctx, tknKey)
		return nil, errors.Join(err2, ErrTokenExpired)
	}
	obj, err := o.objects.Get(ctx, tkn.Bucket, tkn.Filename)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (o *Storage) Delete(ctx context.Context, userID, filename string) error {
	if err := o.objects.Delete(ctx, userID, filename); err != nil {
		return err
	}
	o.tokens.Delete(ctx, filename)
	return nil
}

func (o *Storage) List(ctx context.Context, bucket string, limit int, lastKey string) *Iterator {
	iter := o.objects.ObjectIterator(ctx, bucket, limit, lastKey)
	return iter
}

/// ----
