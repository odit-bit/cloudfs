package blob

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

func (o *Blobs) Download(ctx context.Context, userID, filename string) (*ObjectInfo, error) {
	return o.objects.Get(ctx, userID, filename)
}

func (o *Blobs) Object(ctx context.Context, userID, filename string) (*ObjectInfo, error) {
	return o.objects.Get(ctx, userID, filename)
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

// type UploadParam struct {
// 	Bucket      string
// 	Filename    string
// 	Size        int64
// 	Body        io.Reader
// 	ContentType string
// }

// func (o *Blobs) Upload(ctx context.Context, param UploadParam) (*ObjectInfo, error) {
// 	return o.objects.Put(ctx, param.Bucket, param.Filename, param.Body, param.Size, param.ContentType)
// }

// type OptionUpload struct {
// 	ContentType string
// }
// type UploadResponse struct {
// }

// type ChunkUploader struct {
// 	ctx     context.Context
// 	total   int64
// 	written int64
// 	errC    chan error
// 	resultC chan *ObjectInfo
// 	result  *ObjectInfo
// 	pw      io.WriteCloser
// 	blobs   *Blobs
// }

// func (c *ChunkUploader) Write(p []byte) (int, error) {
// 	if c.total <= int64((len(p) + int(c.written))) {
// 		return 0, fmt.Errorf("written bytes is more than total size")
// 	}
// 	n, err := c.pw.Write(p)
// 	if err != nil {
// 		return n, err
// 	}
// 	c.written += int64(n)
// 	return n, nil
// }

// func (c *ChunkUploader) Result() *ObjectInfo {
// 	return c.result
// }

// func (c *ChunkUploader) Close() error {
// 	err := c.pw.Close()
// 	select {
// 	case err1 := <-c.errC:
// 		err = errors.Join(err1)
// 	case <-c.ctx.Done():
// 		err = errors.Join(err, c.ctx.Err())
// 	}

// 	res, ok := <-c.resultC
// 	if !ok {
// 		return errors.Join(fmt.Errorf("result is nil"), err)
// 	}
// 	c.result = res

// 	return err
// }

// func (o *Blobs) StartChunkUpload(ctx context.Context, bucket, filename string, size int64, opt OptionUpload) *ChunkUploader {
// 	pr, pw := io.Pipe()
// 	errC := make(chan error, 1)
// 	resC := make(chan *ObjectInfo, 1)

// 	go func(err chan<- error) {
// 		info, inErr := o.Upload(ctx, UploadParam{
// 			Bucket:      bucket,
// 			Filename:    filename,
// 			Size:        size,
// 			ContentType: opt.ContentType,
// 			Body:        pr,
// 		})
// 		if inErr != nil {
// 			err <- inErr
// 		} else {
// 			resC <- info
// 		}
// 		close(errC)
// 		close(resC)
// 	}(errC)

// 	return &ChunkUploader{
// 		ctx:     ctx,
// 		total:   size,
// 		written: 0,
// 		pw:      pw,
// 		blobs:   o,
// 		errC:    errC,
// 		resultC: resC,
// 	}
// }

func (o *Blobs) CreateShareToken(ctx context.Context, userID, filename string, expire time.Duration) (*Token, error) {
	tkn := NewShareToken(userID, filename, expire)
	if err := o.tokens.Put(ctx, tkn); err != nil {
		return tkn, err
	}
	return tkn, nil
}

func (o *Blobs) DownloadToken(ctx context.Context, tknKey string) (*ObjectInfo, error) {
	tkn, ok, err := o.tokens.Get(ctx, tknKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	if !tkn.IsNotExpire() {
		err2 := o.tokens.Delete(ctx, tknKey)
		return nil, errors.Join(err2, fmt.Errorf("token expired"))
	}
	return o.objects.Get(ctx, tkn.UserID(), tkn.Filename())
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
