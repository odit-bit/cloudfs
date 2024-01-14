package blob

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const DefaultEndpoint = "127.0.0.1:9000"

var riPool = sync.Pool{
	New: func() any {
		var ri RetreiveInfo
		return &ri
	},
}

// object represent
type RetreiveInfo struct {
	UserID       string
	ObjName      string
	ContentType  string
	Sum          string
	Size         int64
	LastModified time.Time
	Reader       io.ReadCloser `json:"_"`
}

func (ri *RetreiveInfo) Reset() {
	ri.UserID = ""
	ri.ObjName = ""
	ri.ContentType = ""
	ri.Sum = ""
	ri.Size = 0
	ri.LastModified = time.Time{}
	ri.Reader = nil
	riPool.Put(ri)
}

// var prPool = sync.Pool{
// 	New: func() any {
// 		var pr putRequest
// 		return &pr
// 	},
// }

// func NewPutRequest() *putRequest {
// 	pr := prPool.Get().(*putRequest)
// 	return pr
// }

// type putRequest struct {
// 	UserID      string
// 	ObjName     string
// 	ContentType string
// 	Size        int64
// 	Reader      io.Reader
// }

// func (pr *putRequest) Close() {
// 	pr.UserID = ""
// 	pr.ObjName = ""
// 	pr.ContentType = ""
// 	pr.Size = 0
// 	pr.Reader = nil

// 	prPool.Put(pr)
// }

type Result struct {
	Sum       string
	Timestamp time.Time
}

// represent blob storage
type Storage struct {
	minioCli *minio.Client
}

func NewStorage(addr, key, secret string) (*Storage, error) {
	endpoint := addr          //"localhost:9000"
	accessKeyID := key        //"admin"          //"i4mZYsUZpPG9bFXwWRMI"
	secretAccessKey := secret //"admin12345" //"JyNWoAsAzf7KCdSiqTHfzMWtix862PGvnTpKeYCp"
	secure := false

	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	cancel, err := cli.HealthCheck(2 * time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()

	if ok := cli.IsOnline(); !ok {
		return nil, fmt.Errorf("storage endpoint is offline, api-endpoint: %v", cli.EndpointURL().String())
	}

	bs := Storage{
		minioCli: cli,
	}
	return &bs, nil
}

// check if bucket exist, if not it will create new One
func (bs *Storage) MakeBucket(ctx context.Context, userID string) error {
	ok, err := bs.minioCli.BucketExists((ctx), userID)
	if err != nil {
		return err
	}

	if !ok {
		err := bs.minioCli.MakeBucket(ctx, userID, minio.MakeBucketOptions{
			Region: "arab-selatan",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (bs *Storage) Retreive(ctx context.Context, userID string, objName string) (*RetreiveInfo, error) {

	obj, err := bs.minioCli.GetObject(ctx, userID, objName, minio.GetObjectOptions{
		ServerSideEncryption: nil,
		VersionID:            "",
		PartNumber:           0,
		Checksum:             false,
		Internal:             minio.AdvancedGetOptions{},
	})
	if err != nil {
		return nil, err
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	ri := RetreiveInfo{
		UserID:       userID,
		ObjName:      objName,
		ContentType:  stat.ContentType,
		Sum:          stat.ETag,
		Size:         stat.Size,
		LastModified: stat.LastModified,
		Reader:       obj,
	}

	return &ri, nil
}

// save object from pr , if success return sum
func (bs *Storage) Save(ctx context.Context, userID, ObjName, contentType string, size int64, r io.Reader) (string, error) {
	info, err := bs.minioCli.PutObject(ctx, userID, ObjName, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return info.ETag, nil
}

type ListConfig struct {
	UserID string
	Marker string
	Max    int
}

// create a pagination-like request-response
func (bs *Storage) List(ctx context.Context, userID string, max int, startFrom string) StoreIterator {

	c := bs.minioCli.ListObjects(ctx, userID, minio.ListObjectsOptions{
		WithVersions: false,
		WithMetadata: false,
		Prefix:       "",
		Recursive:    false,
		MaxKeys:      max,
		StartAfter:   "",
		UseV1:        false,
	})

	iter := Iterator{
		UserID:  userID,
		listC:   c,
		objInfo: &minio.ObjectInfo{},
		err:     nil,
	}

	return &iter

}

type Iterator struct {
	UserID  string
	listC   <-chan minio.ObjectInfo
	objInfo *minio.ObjectInfo

	err error
}

func (li *Iterator) Next() bool {
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

func (li *Iterator) Value() *RetreiveInfo {

	ri := getRetrieveInfo()
	ri.UserID = li.UserID
	ri.ObjName = li.objInfo.Key
	ri.ContentType = li.objInfo.ContentType
	ri.Sum = li.objInfo.ETag
	ri.Size = li.objInfo.Size
	ri.LastModified = li.objInfo.LastModified

	li.objInfo = nil
	return ri
}

func (li *Iterator) Error() error {
	return li.err
}

func getRetrieveInfo() *RetreiveInfo {
	ri := riPool.Get().(*RetreiveInfo)
	return ri
}

// ====================

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
