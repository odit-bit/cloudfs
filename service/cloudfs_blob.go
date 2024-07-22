package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/odit-bit/cloudfs/internal/blob"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrBucketNotExisted   = errors.New("bucket not existed")
	ErrBucketExisted      = errors.New("bucket already exist")
	ErrAccountExist       = errors.New("account already exist")
	ErrTokenExpired       = errors.New("invalid token")
	ErrFileNotExisted     = errors.New("file not found")

	ErrUpload = errors.New("upload error")
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
	// if param.Size <= 0 {
	// 	return &errUpload{msg: "size cannot 0 or under zero"}
	// }
	// if param.ContentType == ""{}
	if param.DataReader == nil {
		return &errUpload{msg: "param data reader is nil"}
	}

	return nil

}

func (app *Cloudfs) Upload(ctx context.Context, param *UploadParam) (*blob.ObjectInfo, error) {
	if err := param.validate(); err != nil {
		return nil, err
	}

	obj, err := app.blobService.Put(
		ctx,
		param.UserID,
		param.Filename,
		param.DataReader,
		param.Size,
		param.ContentType,
	)
	if err != nil {
		return nil, err
	}

	return obj, err
}

type DownloadParam struct {
	UserID   string
	Filename string
}

type sharingObj struct {
	Owner      string
	Filename   string
	Token      string
	ValidUntil string
}

func (app *Cloudfs) SharingFile(ctx context.Context, userID, filename string) (*sharingObj, error) {
	_, err := app.blobService.Get(ctx, userID, filename)
	if err != nil {
		return nil, err
	}

	// token, err := app.tokenService.Generate(ctx, userID, filename, 24*time.Hour)
	// if err != nil {
	// 	return nil, err
	// }

	// sObj := sharingObj{
	// 	Owner:      userID,
	// 	Filename:   filename,
	// 	Token:      token,
	// 	ValidUntil: humanize.Time(time.Now().Add(25 * time.Hour)),
	// }

	// return &sObj, nil

	var sObj sharingObj
	if err := app.tokenService.Query(ctx, func(txn TokenTxn) error {
		tkn := blob.NewShareToken(userID, filename, 24*time.Hour)
		if err := txn.Put(ctx, tkn); err != nil {
			txn.Cancel()
			return err
		}

		sObj.Owner = userID
		sObj.Filename = filename
		sObj.Token = tkn.Key
		sObj.ValidUntil = tkn.Expire.String()

		return txn.Commit()
	}); err != nil {
		return nil, err
	}
	return &sObj, nil
}

func (app *Cloudfs) DownloadSharedFile(ctx context.Context, token string, writeFunc func(r io.Reader)) error {
	// userID, filename, ok := app.tokenService.Validate(ctx, token)
	// if !ok {
	// 	return ErrTokenExpired
	// }

	var userID, filename string
	if err := app.tokenService.Query(ctx, func(txn TokenTxn) error {
		st, err := txn.Get(ctx, token)
		if err != nil {
			txn.Cancel()
			return fmt.Errorf("token not Found")
		}

		if ok := st.IsNotExpire(); !ok {
			if err := txn.Delete(ctx, st.Key); err != nil {
				txn.Cancel()
				return ErrTokenExpired
			}
		}

		userID = st.UserID
		filename = st.Filename
		return txn.Commit()
	}); err != nil {
		return err
	}

	info, err := app.blobService.Get(ctx, userID, filename)
	if err != nil {
		return ErrFileNotExisted
	}

	reader, err := info.Reader.Get(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	writeFunc(reader)
	return nil

}

func (app *Cloudfs) Object(ctx context.Context, userID, filename string) (*blob.ObjectInfo, error) {
	return app.blobService.Get(ctx, userID, filename)
}

func (app *Cloudfs) Delete(ctx context.Context, userID, filename string) error {
	return app.blobService.Delete(ctx, userID, filename)
}

func (app *Cloudfs) Download(ctx context.Context, w io.Writer, param *DownloadParam) (any, error) {
	// if ok, err := app.bucketService.IsBucketExist(ctx, param.UserID); err != nil {
	// 	return nil, err
	// } else if !ok {
	// 	return nil, ErrBucketNotExisted
	// }

	obj, err := app.blobService.Get(ctx, param.UserID, param.Filename)
	if err != nil {
		return nil, err
	}

	src, err := obj.Reader.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	_, err = io.Copy(w, src)

	return nil, err
}

func (app *Cloudfs) ListObject(ctx context.Context, bucket string, limit int, lastKey string) ([]*blob.ObjectInfo, error) {
	iter := app.blobService.ObjectIterator(ctx, bucket, limit, lastKey)
	if iter.Error() != nil {
		return nil, iter.Error()
	}
	infos := make([]*blob.ObjectInfo, limit)

	count := 0
	for iter.Next() {
		infos[count] = iter.Value()
		count++
	}

	return infos[:count], iter.Error()
}
