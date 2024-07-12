package blob

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const DefaultEndpoint = "127.0.0.1:9000"

// object represent
type ObjectInfo struct {
	UserID       string
	ObjName      string
	ContentType  string
	Sum          string
	Size         int64
	LastModified time.Time
	reader       io.ReadCloser
}

func (oi *ObjectInfo) Reader() io.Reader {
	return oi.reader
}

func (oi *ObjectInfo) Close() error {
	return oi.reader.Close()
}

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
type Storage struct {
	minioCli *minio.Client
}

func NewStorage(addr, key, secret string) (*Storage, error) {
	// 	endpoint := addr          //"localhost:9000"
	// 	accessKeyID := key        //"admin"          //"i4mZYsUZpPG9bFXwWRMI"
	// 	secretAccessKey := secret //"admin12345" //"JyNWoAsAzf7KCdSiqTHfzMWtix862PGvnTpKeYCp"
	// secure := false

	minioCli, err := connectMinio(addr, key, secret, false)
	if err != nil {
		return nil, err
	}

	bs := Storage{
		minioCli: minioCli,
	}
	return &bs, nil
}

func (s *Storage) CreateBucket(ctx context.Context, bucketName string) error {
	ok, err := s.minioCli.BucketExists((ctx), bucketName)
	if err != nil {
		return err
	}

	if !ok {
		err := s.minioCli.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: "arab-selatan",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) Put(ctx context.Context, userID, filename string, file io.Reader, size int64, contentType string) (any, error) {
	return s.minioCli.PutObject(ctx, userID, filename, file, size, minio.PutObjectOptions{})
}

func (s *Storage) Get(ctx context.Context, bucketName, filename string, objInfo *ObjectInfo) error {
	res, err := s.minioCli.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer res.Close()

	stat, _ := res.Stat()
	objInfo.UserID = stat.Owner.DisplayName
	objInfo.ObjName = stat.Key
	objInfo.ContentType = stat.ContentType
	objInfo.Sum = stat.ChecksumSHA256
	objInfo.Size = stat.Size
	objInfo.LastModified = stat.LastModified
	objInfo.reader = res

	return nil
}

func (s *Storage) List(ctx context.Context, bucketName string, limit int, lastKey string) *Cursor {
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

type Cursor struct {
	UserID  string
	listC   <-chan minio.ObjectInfo
	objInfo *minio.ObjectInfo

	err error
}

func (li *Cursor) Next() bool {
	objInfo, ok := <-li.listC
	if !ok {
		return false
	}
	if objInfo.Err != nil {
		li.err = objInfo.Err
		return false
	}

	li.objInfo = &objInfo
	return true
}

func (li *Cursor) Scan(info *ObjectInfo) {
	if info == nil {
		panic("storage cursor cannot scan into nil info")
	}

	// fmt.Println("STORAGE CURSOR", li.objInfo)
	info.UserID = li.UserID
	info.ObjName = li.objInfo.Key
	info.Sum = li.objInfo.ETag
	info.Size = li.objInfo.Size
	info.LastModified = li.objInfo.LastModified

	li.objInfo = nil
}

func (li *Cursor) Error() error {
	return li.err
}

// ====================

func parseUploadRequest(userID string, minioClient *minio.Client, r *http.Request) (any, error) {
	// size, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	// if err != nil {
	// 	return nil, err
	// }

	reader, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	part, err := reader.NextPart()
	if err != nil {
		return nil, err
	}
	defer part.Close()

	if part.FormName() != "file" {
		return nil, err
	}

	// save blob into storage
	entry := minio.PutObjectFanOutEntry{
		Key:         part.FileName(),
		ContentType: part.Header.Get("content-type"),
	}

	_, err = minioClient.PutObjectFanOut(r.Context(), userID, part, minio.PutObjectFanOutRequest{
		Entries: []minio.PutObjectFanOutEntry{entry},
	})

	return nil, err

}

// for testing

func (bs *Storage) purge(userID string) error {
	objsC := bs.minioCli.ListObjects(context.Background(), userID, minio.ListObjectsOptions{})
	errC := bs.minioCli.RemoveObjects(context.Background(), userID, objsC, minio.RemoveObjectsOptions{})
	for objErr := range errC {
		if objErr.Err != nil {
			return objErr.Err
		}
	}
	return nil
}
