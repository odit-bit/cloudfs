package repo

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/odit-bit/cloudfs/internal/blob"
)

// const DefaultEndpoint = "127.0.0.1:9000"

// func connectMinio(endpoint, accessKeyID, secretAccessKey string, secure bool) (*minio.Client, error) {
// 	// fmt.Println(endpoint, accessKeyID, secretAccessKey)
// 	cli, err := minio.New(endpoint, &minio.Options{
// 		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
// 		Secure: secure,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	cancel, err := cli.HealthCheck(5 * time.Second)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cancel()

// 	if ok := cli.IsOnline(); !ok {
// 		return nil, fmt.Errorf("storage endpoint is offline, api-endpoint: %v", cli.EndpointURL().String())
// 	}

// 	return cli, nil

// }

// func NewMinioBlob(addr, key, secret string) (*MinioAdapter, error) {
// 	// 	endpoint := addr
// 	// 	accessKeyID := key
// 	// 	secretAccessKey := secret
// 	// secure := false

// 	minioCli, err := connectMinio(addr, key, secret, false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	bs := MinioAdapter{
// 		minioCli: minioCli,
// 	}
// 	return &bs, nil
// }

// represent blob storage
type MinioAdapter struct {
	minioCli *minio.Client
}

func NewMinioBlob(cli *minio.Client) *MinioAdapter {
	return &MinioAdapter{minioCli: cli}
}

func (s *MinioAdapter) CreateBucket(ctx context.Context, bucketName string) (any, error) {
	ok, err := s.minioCli.BucketExists((ctx), bucketName)
	if err != nil {
		return nil, err
	}

	if !ok {
		err := s.minioCli.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (s *MinioAdapter) Delete(ctx context.Context, bucketName, key string) error {
	return s.minioCli.RemoveObject(ctx, bucketName, key, minio.RemoveObjectOptions{})
}

func (s *MinioAdapter) Put(ctx context.Context, bucketName, key string, file io.ReadCloser, size int64, contentType string) (*blob.Info, error) {
	info, err := s.put(ctx, bucketName, key, file, size, contentType)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *MinioAdapter) put(ctx context.Context, bucketName, key string, file io.ReadCloser, size int64, contentType string) (*blob.Info, error) {
	defer file.Close()

	if ok, err := s.minioCli.BucketExists(ctx, bucketName); !ok {
		//create bucket
		if err := s.minioCli.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	res, err := s.minioCli.PutObject(ctx, bucketName, key, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	info := blob.Info{
		Bucket:       bucketName,
		Filename:     key,
		ContentType:  contentType,
		Sum:          res.ChecksumSHA256,
		Size:         res.Size,
		LastModified: res.LastModified,
		// Reader:       nil,
	}
	return &info, nil
}

func (s *MinioAdapter) GetUsage(ctx context.Context, bucket string) (int64, error) {
	return 0, fmt.Errorf("quota is not implemented")
}

func (s *MinioAdapter) Get(ctx context.Context, bucketName, filename string) (*blob.Object, error) {
	// stat, err := s.minioCli.StatObject(ctx, bucketName, filename, minio.GetObjectOptions{})

	data, err := s.minioCli.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	stat, err := data.Stat()
	if err != nil {
		return nil, err
	}

	var objInfo blob.Object
	objInfo.Bucket = bucketName
	objInfo.Filename = stat.Key
	objInfo.ContentType = stat.ContentType
	objInfo.Sum = stat.ChecksumSHA256
	objInfo.Size = stat.Size
	objInfo.LastModified = stat.LastModified
	objInfo.Data = data
	return &objInfo, nil
}

func (s *MinioAdapter) Info(ctx context.Context, bucketName, fileName string) (*blob.Info, error) {
	stat, err := s.minioCli.StatObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	var objInfo blob.Info
	objInfo.Bucket = bucketName
	objInfo.Filename = stat.Key
	objInfo.ContentType = stat.ContentType
	objInfo.Sum = stat.ChecksumSHA256
	objInfo.Size = stat.Size
	objInfo.LastModified = stat.LastModified
	return &objInfo, nil
}

func (s *MinioAdapter) GetShareUrl(bucket, key string) blob.ShareFunc {
	return func(ctx context.Context, expiration time.Duration) (*url.URL, error) {
		return s.minioCli.PresignedGetObject(ctx, bucket, key, expiration, url.Values{})
	}
}

func (s *MinioAdapter) ObjectIterator(ctx context.Context, bucketName string, limit int, lastFilename string) *blob.Iterator {
	res := s.minioCli.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:     "",
		MaxKeys:    limit,
		StartAfter: lastFilename,
	})

	objC := make(chan *blob.Info)

	go func() {
		for c := range res {
			obj := &blob.Info{}
			//NOTE: return minio.ObjectInfo only 4 field that has value,
			//see https://min.io/docs/minio/linux/developers/go/API.html#ListObjects
			obj.Filename = c.Key
			obj.Sum = c.ETag
			obj.Size = c.Size
			obj.ContentType = c.ContentType // always ""
			obj.LastModified = c.LastModified

			// the last object from minio will always contain Err
			if c.Err == nil {
				objC <- obj
			}
		}
		close(objC)
	}()

	iter := blob.Iterator{
		UserID: "",
		C:      objC,
	}
	return &iter
}

// func (s *MinioAdapter) MakeBucket(ctx context.Context, bucketName, region string) error {
// 	return s.minioCli.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
// }

// func (s *MinioAdapter) IsBucketExist(ctx context.Context, bucketName string) (bool, error) {
// 	return s.minioCli.BucketExists(ctx, bucketName)
// }

// ====================

// for testing

func (bs *MinioAdapter) sharedURL(bucketName, key string, dur time.Duration) {
}

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
