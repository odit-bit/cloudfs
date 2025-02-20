package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

type ObjectStorer interface {
	Put(ctx context.Context, bucket, filename string, reader io.Reader, size int64, contentType string) (*ObjectInfo, error)
	Get(ctx context.Context, bucket, filename string) (*ObjectInfo, error)
	Delete(ctx context.Context, bucket, filename string) error
	ObjectIterator(ctx context.Context, bucket string, limit int, lastFilename string) Iterator
}

type TokenStorer interface {
	Put(ctx context.Context, token *Token) error
	Get(ctx context.Context, tokenKey string) (*Token, bool, error)
	Delete(ctx context.Context, tokenKey string) error
}

type AccountStorer interface {
	Authorize(ctx context.Context)
}

type Blobs struct {
	tokens  TokenStorer
	objects ObjectStorer
}

func NewWithMemory() (*Blobs, error) {
	blobs, _ := newObjectMemory()
	tokens, _ := newObjectTokenMemory()
	return NewBlobs(context.Background(), tokens, blobs)
}

func NewBlobs(ctx context.Context, tokens TokenStorer, objects ObjectStorer) (*Blobs, error) {
	return &Blobs{tokens: tokens, objects: objects}, nil
}

func (o *Blobs) Download(ctx context.Context, bucket, filename string) (*ObjectInfo, error) {
	return o.objects.Get(ctx, bucket, filename)
}

func (o *Blobs) Object(ctx context.Context, bucket, filename string) (*ObjectInfo, error) {
	return o.objects.Get(ctx, bucket, filename)
}

type UploadParam struct {
	Bucket      string
	Filename    string
	Size        int64
	Body        io.Reader
	ContentType string
}

func (o *Blobs) Upload(ctx context.Context, param UploadParam) (*ObjectInfo, error) {
	return o.objects.Put(ctx, param.Bucket, param.Filename, param.Body, param.Size, param.ContentType)
}

func (o *Blobs) CreateShareToken(ctx context.Context, userID, filename string, expire time.Duration) (*Token, error) {
	tkn := NewShareToken(userID, filename, expire)
	if err := o.tokens.Put(ctx, tkn); err != nil {
		return tkn, err
	}
	return tkn, nil
}

var (
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidShareToken = errors.New("invalid share token")
	ErrUnknown           = errors.New("exception")
)

func (o *Blobs) DownloadToken(ctx context.Context, tknKey string) (*ObjectInfo, error) {
	tkn, ok, err := o.tokens.Get(ctx, tknKey)
	if err != nil {
		return nil, errors.Join(ErrUnknown, err)
	}
	if !ok {
		return nil, ErrInvalidShareToken
	}
	if !tkn.IsNotExpire() {
		err2 := o.tokens.Delete(ctx, tknKey)
		return nil, errors.Join(err2, ErrTokenExpired)
	}
	obj, err := o.objects.Get(ctx, tkn.UserID(), tkn.Filename())
	if err != nil {
		return nil, errors.Join(err, ErrUnknown)
	}
	return obj, nil
}

func (o *Blobs) Delete(ctx context.Context, userID, filename string) error {
	return o.objects.Delete(ctx, userID, filename)
}

func (o *Blobs) ListObject(ctx context.Context, bucket string, limit int, lastKey string) ([]ObjectInfo, error) {
	iter := o.objects.ObjectIterator(ctx, bucket, limit, lastKey)
	if iter.Error() != nil {
		return nil, iter.Error()
	}
	infos := make([]ObjectInfo, limit)

	count := 0
	for iter.Next() {
		infos[count] = iter.Value()
		count++
	}

	return infos[:count], iter.Error()
}

/// ----

type ChunkWriter struct {
	ctx     context.Context
	pr      *io.PipeReader
	pw      *io.PipeWriter
	resultC chan any
	errC    chan error
	done    chan struct{}
}

func (o *Blobs) NewChunkWriter(ctx context.Context, bucket, filename, contentType string, size int64) *ChunkWriter {
	pr, pw := io.Pipe()
	resultC := make(chan any, 1)
	errC := make(chan error, 1)
	done := make(chan struct{}, 1)
	go func() {
		defer pr.Close()
		res, err := o.objects.Put(ctx, bucket, filename, pr, size, contentType)
		if err != nil {
			errC <- err
		}
		resultC <- res
		close(errC)
		close(resultC)
		close(done)
	}()
	return &ChunkWriter{
		ctx:     ctx,
		pw:      pw,
		pr:      pr,
		resultC: resultC,
		errC:    errC,
		done:    done,
	}
}

func (cw *ChunkWriter) Write(p []byte) (int, error) {
	return cw.pw.Write(p)
}

// will close writer and return the result, any write call after this will error
func (cw *ChunkWriter) Result() (*ObjectInfo, error) {
	err := cw.pw.Close()
	if err != nil {
		return nil, err
	}

	select {
	case <-cw.ctx.Done():
		return nil, cw.ctx.Err()
	case <-cw.done:
		err, ok := <-cw.errC
		if ok {
			return nil, err
		}
		res, ok := <-cw.resultC
		if !ok {
			return nil, fmt.Errorf("error and result channel is nil, this is a bug")
		}
		info, ok := res.(*ObjectInfo)
		if !ok {
			panic(fmt.Sprintf("invalid type %T", res))
		}
		return info, nil
	}
}
