package blob

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

var _ ObjectStorer = (*aferoBlob)(nil)

type aferoBlob struct {
	Fs afero.Fs
}

func newObjectMemory() (*aferoBlob, error) {
	return &aferoBlob{
		Fs: afero.NewMemMapFs(),
	}, nil

}

func newAferoBlob(root string) (*aferoBlob, error) {
	return &aferoBlob{
		Fs: afero.NewBasePathFs(afero.NewOsFs(), root),
	}, nil
}

func (store *aferoBlob) GetUsage(ctx context.Context, bucket string) (int64, error) {
	return 0, fmt.Errorf("local disk usage not implemented")
}

// ObjectIterator implements service.BlobStore.
func (store *aferoBlob) ObjectIterator(ctx context.Context, bucket string, limit int, lastFilename string) *Iterator {
	dirPath := bucket
	c := make(chan *Info)
	if ok, _ := afero.DirExists(store.Fs, dirPath); !ok {
		close(c)
		return &Iterator{
			C: c,
		}
	}

	entry, err := afero.ReadDir(store.Fs, dirPath)
	if err != nil {
		close(c)
		return &Iterator{
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
		for idx, info := range entry {
			if idx >= limit {
				break
			}

			info := &Info{
				Bucket:       bucket,
				Filename:     info.Name(),
				ContentType:  "",
				Sum:          "",
				Size:         info.Size(),
				LastModified: info.ModTime(),
				// Reader:       nil,
			}

			select {
			case <-ctx.Done():
				close(c)
				return
			case c <- info:
				continue
			}
		}
		close(c)
	}()

	return &Iterator{
		UserID: "",
		C:      c,
	}
}

func (store *aferoBlob) Delete(ctx context.Context, userID, filename string) error {
	fileKey := filepath.Join(userID, filename)
	return store.Fs.Remove(fileKey)
}

func (store *aferoBlob) Get(ctx context.Context, userID string, filename string) (*Object, error) {
	fileKey := filepath.Join(userID, filename)
	ok, err := afero.Exists(store.Fs, fileKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("file not existed")
	}

	f, err := store.Fs.Open(fileKey)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	obj := &Object{
		Info: Info{
			Bucket:       userID,
			Filename:     filename,
			ContentType:  "",
			Sum:          "",
			Size:         stat.Size(),
			LastModified: stat.ModTime(),
		},
		Data: f,
	}

	return obj, nil
}

func (store *aferoBlob) Put(ctx context.Context, userID string, filename string, reader io.ReadCloser, size int64, contentType string) (*Info, error) {
	fileKey := filepath.Join(userID, filename)
	//optimis path
	ok, err := afero.Exists(store.Fs, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		if err := store.Fs.Mkdir(userID, os.ModePerm); err != nil {
			return nil, err
		}
	}

	f, err := store.Fs.Create(fileKey)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hasher := sha256.New()

	tr := io.TeeReader(reader, hasher)
	if _, err := io.Copy(f, tr); err != nil {
		return nil, err
	}

	sum := hex.EncodeToString(hasher.Sum(nil))
	stat, _ := f.Stat()
	return &Info{
		Bucket:       userID,
		Filename:     filename,
		ContentType:  contentType,
		Sum:          sum,
		Size:         size,
		LastModified: stat.ModTime(),
		// Reader:       nil,
	}, nil
}

func Hasher256(r io.Reader) string {
	hasher := sha256.New()
	io.Copy(hasher, r)
	return hex.EncodeToString(hasher.Sum(nil))
}
