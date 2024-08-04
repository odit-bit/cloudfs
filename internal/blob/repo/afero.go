package repo

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/spf13/afero"
)

type aferoBlob struct {
	Fs afero.Fs
}

func NewInMemBlob() (aferoBlob, error) {
	return aferoBlob{
		Fs: afero.NewMemMapFs(),
	}, nil

}

func NewAferoBlob(root string) (*aferoBlob, error) {
	return &aferoBlob{
		Fs: afero.NewOsFs(),
	}, nil
}

// ObjectIterator implements service.BlobStore.
func (store *aferoBlob) ObjectIterator(ctx context.Context, userID string, limit int, lastFilename string) blob.Iterator {
	dirPath := userID
	c := make(chan blob.ObjectInfo)
	if ok, _ := afero.DirExists(store.Fs, dirPath); !ok {
		close(c)
		return blob.Iterator{
			C: c,
		}
	}

	entry, err := afero.ReadDir(store.Fs, dirPath)
	if err != nil {
		close(c)
		return blob.Iterator{
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

			info := blob.ObjectInfo{
				UserID:       userID,
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

	return blob.Iterator{
		UserID: "",
		C:      c,
	}
}

func (store *aferoBlob) Delete(ctx context.Context, userID, filename string) error {
	fileKey := filepath.Join(userID, filename)
	return store.Fs.Remove(fileKey)
}

func (store *aferoBlob) Get(ctx context.Context, userID string, filename string) (*blob.ObjectInfo, error) {
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

	return &blob.ObjectInfo{
		UserID:       userID,
		Filename:     filename,
		ContentType:  "",
		Sum:          "",
		Size:         stat.Size(),
		LastModified: stat.ModTime(),
		Data:         f,
	}, nil
}

func (store *aferoBlob) Put(ctx context.Context, userID string, filename string, reader io.Reader, size int64, contentType string) (*blob.ObjectInfo, error) {
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
	if _, err := io.Copy(f, reader); err != nil {
		return nil, err
	}

	stat, _ := f.Stat()

	return &blob.ObjectInfo{
		UserID:       userID,
		Filename:     filename,
		ContentType:  contentType,
		Sum:          "",
		Size:         size,
		LastModified: stat.ModTime(),
		// Reader:       nil,
	}, nil
}
