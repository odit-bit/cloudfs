package service

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/token"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrBucketNotExisted   = errors.New("bucket not existed")
	ErrBucketExisted      = errors.New("bucket already exist")
	ErrAccountExist       = errors.New("account already exist")
	ErrTokenExpired       = errors.New("token is expired")
	ErrTokenNotExisted    = errors.New("invalid token")
	ErrFileNotExisted     = errors.New("file not found")

	ErrUpload = errors.New("upload error")
)

const (
	_DEFAULT_FILE_SHARE_TTL = 24 * time.Hour
)

type UploadParam struct {
	UserID      string
	Filename    string
	Size        int64
	ContentType string
	DataReader  io.Reader
}

type errUpload struct {
	msg string
}

func (err *errUpload) Error() string {
	return err.msg
}

func (err *errUpload) Is(tErr error) bool {
	return tErr == ErrUpload
}

func (param *UploadParam) validate() error {
	if param.UserID == "" {
		return &errUpload{msg: "userId is nil"}
	}
	if param.Filename == "" {
		return &errUpload{msg: "filename is nil"}
	}

	if param.DataReader == nil {
		return &errUpload{msg: "param data reader is nil"}
	}

	return nil

}

type UploadInfo struct {
}

func (app *Cloudfs) Upload(ctx context.Context, param *UploadParam) (UploadInfo, error) {
	if err := param.validate(); err != nil {
		return UploadInfo{}, err
	}

	_, err := app.blobs.Put(
		ctx,
		param.UserID,
		param.Filename,
		param.DataReader,
		param.Size,
		param.ContentType,
	)
	if err != nil {
		return UploadInfo{}, err
	}

	return UploadInfo{}, err
}

type DownloadParam struct {
	UserID   string
	Filename string
}

func (app *Cloudfs) SharingFile(ctx context.Context, userID, filename string) (*token.ShareToken, error) {
	_, err := app.blobs.Get(ctx, userID, filename)
	if err != nil {
		return nil, err
	}

	tkn := token.NewShareToken(userID, filename, _DEFAULT_FILE_SHARE_TTL)

	if err := app.tokens.Put(ctx, tkn); err != nil {
		return nil, err
	}

	return tkn, nil
}

func (app *Cloudfs) DownloadSharedFile(ctx context.Context, token string) (*blob.ObjectInfo, error) {

	tkn, ok, err := app.tokens.Get(ctx, token)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrTokenNotExisted
	}

	if !tkn.IsNotExpire() {
		return nil, ErrTokenExpired
	}

	info, err := app.blobs.Get(ctx, tkn.UserID(), tkn.Filename())
	if err != nil {
		return nil, ErrFileNotExisted
	}
	return info, nil

}

func (app *Cloudfs) Object(ctx context.Context, userID, filename string) (*blob.ObjectInfo, error) {
	return app.blobs.Get(ctx, userID, filename)
}

func (app *Cloudfs) Delete(ctx context.Context, userID, filename string) error {
	return app.blobs.Delete(ctx, userID, filename)
}

func (app *Cloudfs) Download(ctx context.Context, w io.Writer, param *DownloadParam) (any, error) {

	obj, err := app.blobs.Get(ctx, param.UserID, param.Filename)
	if err != nil {
		return nil, err
	}
	defer obj.Data.Close()
	_, err = io.Copy(w, obj.Data)

	return nil, err
}

func (app *Cloudfs) ListObject(ctx context.Context, bucket string, limit int, lastKey string) ([]blob.ObjectInfo, error) {
	iter := app.blobs.ObjectIterator(ctx, bucket, limit, lastKey)
	if iter.Error() != nil {
		return nil, iter.Error()
	}
	infos := make([]blob.ObjectInfo, limit)

	count := 0
	for iter.Next() {
		infos[count] = iter.Value()
		count++
	}

	return infos[:count], iter.Error()
}
