package blob

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const DefaultEndpoint = "127.0.0.1:9000"

type Result struct {
	Sum       string
	Timestamp time.Time
}

func connectMinio(endpoint, accessKeyID, secretAccessKey string, secure bool) (*minio.Client, error) {
	fmt.Println(endpoint, accessKeyID, secretAccessKey)
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	cancel, err := cli.HealthCheck(5 * time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()

	if ok := cli.IsOnline(); !ok {
		return nil, fmt.Errorf("storage endpoint is offline, api-endpoint: %v", cli.EndpointURL().String())
	}

	// if err := cli.MakeBucket(context.TODO(), "init-bucket", minio.MakeBucketOptions{}); err != nil {
	// 	panic(err)
	// }

	return cli, nil

}

// represent blob storage
type MinioAdapter struct {
	minioCli *minio.Client
}

func NewMinioAdapter(addr, key, secret string) (*MinioAdapter, error) {
	// 	endpoint := addr          //"localhost:9000"
	// 	accessKeyID := key        //"admin"
	// 	secretAccessKey := secret //"admin12345"
	// secure := false

	minioCli, err := connectMinio(addr, key, secret, false)
	if err != nil {
		return nil, err
	}
	bs := MinioAdapter{
		minioCli: minioCli,
	}
	return &bs, nil
}

func (s *MinioAdapter) CreateBucket(ctx context.Context, bucketName string) (any, error) {
	ok, err := s.minioCli.BucketExists((ctx), bucketName)
	if err != nil {
		return nil, err
	}

	if !ok {
		err := s.minioCli.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: "arab-selatan",
		})
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (s *MinioAdapter) Delete(ctx context.Context, bucketName, key string) error {
	return s.minioCli.RemoveObject(ctx, bucketName, key, minio.RemoveObjectOptions{})
}

func (s *MinioAdapter) Put(ctx context.Context, bucketName, key string, file io.Reader, size int64, contentType string) (*ObjectInfo, error) {
	return s.put(ctx, bucketName, key, file, -1, contentType)
}

func (s *MinioAdapter) put(ctx context.Context, bucketName, key string, file io.Reader, size int64, contentType string) (*ObjectInfo, error) {
	info, err := s.minioCli.PutObject(ctx, bucketName, key, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	obj := ObjectInfo{
		UserID:       bucketName,
		Filename:     key,
		ContentType:  contentType,
		Sum:          info.ETag,
		Size:         info.Size,
		LastModified: info.LastModified,
		// Reader:       nil,
	}
	return &obj, nil
}

func (s *MinioAdapter) Get(ctx context.Context, bucketName, filename string) (*ObjectInfo, error) {
	stat, err := s.minioCli.StatObject(ctx, bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	var objInfo ObjectInfo
	objInfo.UserID = bucketName
	objInfo.Filename = stat.Key
	objInfo.ContentType = stat.ContentType
	objInfo.Sum = stat.ChecksumSHA256
	objInfo.Size = stat.Size
	objInfo.LastModified = stat.LastModified

	objInfo.Reader = s.GetObject(bucketName, filename)

	return &objInfo, nil
}

func (s *MinioAdapter) Info(ctx context.Context, bucketName, fileName string) (*ObjectInfo, error) {
	stat, err := s.minioCli.StatObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	var objInfo ObjectInfo
	objInfo.UserID = bucketName
	objInfo.Filename = stat.Key
	objInfo.ContentType = stat.ContentType
	objInfo.Sum = stat.ChecksumSHA256
	objInfo.Size = stat.Size
	objInfo.LastModified = stat.LastModified
	return &objInfo, nil
}

func (s *MinioAdapter) GetShareUrl(bucket, key string) ShareFunc {
	return func(ctx context.Context, expiration time.Duration) (*url.URL, error) {
		return s.minioCli.PresignedGetObject(ctx, bucket, key, expiration, url.Values{})
	}
}

func (s *MinioAdapter) GetObject(bucket, key string) ReaderFunc {
	return func(ctx context.Context) (io.ReadCloser, error) {
		obj, err := s.minioCli.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
		return obj, err
	}
}

func (s *MinioAdapter) ObjectIterator(ctx context.Context, bucketName string, limit int, lastFilename string) *Iterator {
	res := s.minioCli.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:     "",
		MaxKeys:    limit,
		StartAfter: lastFilename,
	})

	objC := make(chan *ObjectInfo)

	go func() {
		for c := range res {
			var obj ObjectInfo
			//NOTE: return minio.ObjectInfo only 4 field that has value,
			//see https://min.io/docs/minio/linux/developers/go/API.html#ListObjects
			obj.Filename = c.Key
			obj.Sum = c.ETag
			obj.Size = c.Size
			obj.ContentType = c.ContentType // always ""
			obj.LastModified = c.LastModified

			// the last object from minio will always contain Err
			if c.Err == nil {
				objC <- &obj
			}
		}
		close(objC)
	}()

	iter := Iterator{
		UserID: "",
		C:      objC,
		obj:    &ObjectInfo{},
		err:    nil,
	}
	return &iter
}

func (s *MinioAdapter) List(ctx context.Context, bucketName string, limit int, lastKey string) *Cursor {
	res := s.minioCli.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:     "",
		MaxKeys:    limit,
		StartAfter: lastKey,
	})

	return &Cursor{
		UserID:  bucketName,
		listC:   res,
		objInfo: &minio.ObjectInfo{},
		err:     nil,
	}
}

func (s *MinioAdapter) MakeBucket(ctx context.Context, bucketName, region string) error {
	return s.minioCli.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}

func (s *MinioAdapter) IsBucketExist(ctx context.Context, bucketName string) (bool, error) {
	return s.minioCli.BucketExists(ctx, bucketName)
}

// ====================

// for testing

func (bs *MinioAdapter) purge(userID string) error {
	objsC := bs.minioCli.ListObjects(context.Background(), userID, minio.ListObjectsOptions{})
	errC := bs.minioCli.RemoveObjects(context.Background(), userID, objsC, minio.RemoveObjectsOptions{})
	for objErr := range errC {
		if objErr.Err != nil {
			return objErr.Err
		}
	}
	return nil
}
