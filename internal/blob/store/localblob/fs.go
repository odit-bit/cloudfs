package localblob

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/service"
)

var _ service.BlobStore = (*Volume)(nil)

type Volume struct {
	root string
}

func New(path string) (*Volume, error) {
	// it should create dir if not exist
	if err := os.MkdirAll(path, 0o777); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path, 0o777); err != nil {
				return nil, err
			}
		}
	}

	v := Volume{
		root: path,
	}
	return &v, nil
}

func (v *Volume) purge() error {
	return os.RemoveAll(v.root)
}

// Delete implements service.BlobStore.
func (v *Volume) Delete(ctx context.Context, bucket string, filename string) error {
	path := filepath.Join(v.root, bucket, filename)
	return os.RemoveAll(path)
}

// Get implements service.BlobStore.
func (v *Volume) Get(ctx context.Context, bucket string, filename string) (*blob.ObjectInfo, error) {
	path := filepath.Join(v.root, bucket, filename)
	// if ok := checkFile(path); !ok {
	// 	return nil, fmt.Errorf("localblob: file not exist, path %s", path)
	// }

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, _ := os.Open(path)

	readFunc := blob.ReaderFunc(func(ctx context.Context) (io.ReadCloser, error) {
		return f, nil
	})

	obj := blob.ObjectInfo{
		UserID:       bucket,
		Filename:     filename,
		ContentType:  "",
		Sum:          "",
		Size:         stat.Size(),
		LastModified: stat.ModTime(),
		Reader:       readFunc,
	}

	return &obj, nil
}

// ObjectIterator implements service.BlobStore.
func (v *Volume) ObjectIterator(ctx context.Context, bucket string, limit int, lastFilename string) *blob.Iterator {
	dirPath := filepath.Join(v.root, bucket)
	c := make(chan *blob.ObjectInfo)
	if ok := checkDir(dirPath); !ok {
		close(c)
		return &blob.Iterator{
			C: c,
		}
	}

	entry, err := os.ReadDir(dirPath)
	if err != nil {
		close(c)
		return &blob.Iterator{
			UserID: "",
			C:      c,
		}
	}

	// forward the entry until find lastFilename
	// fileDir := filepath.Join(dirPath, lastFilename)
	if lastFilename != "" {
		for n := 0; n+1 < len(entry); n++ {
			v := entry[n]
			if lastFilename == v.Name() {
				entry = entry[n+1:]
				break
			}
		}
	}
	go func() {
		for i, e := range entry {
			if i >= limit {
				break
			}
			info, err := e.Info()
			if err != nil {
				break
			}

			c <- &blob.ObjectInfo{
				UserID:       bucket,
				Filename:     info.Name(),
				ContentType:  "",
				Sum:          "",
				Size:         info.Size(),
				LastModified: info.ModTime(),
				Reader:       nil,
			}
		}
		close(c)
	}()
	return &blob.Iterator{
		UserID: "",
		C:      c,
	}
}

// Put implements service.BlobStore.
func (v *Volume) Put(ctx context.Context, bucket string, filename string, reader io.Reader, size int64, contentType string) (*blob.ObjectInfo, error) {
	// check dir is exist
	path := filepath.Join(v.root, bucket)
	if ok := checkDir(path); !ok {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	f, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return nil, err
	}

	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return nil, err
	}

	return &blob.ObjectInfo{
		UserID:       bucket,
		Filename:     filename,
		ContentType:  contentType,
		Sum:          "",
		Size:         size,
		LastModified: time.Now(),
		Reader:       nil,
	}, nil
}

func checkDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}

	return f.IsDir()
}

func checkFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !f.IsDir()
}
