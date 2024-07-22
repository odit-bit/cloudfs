package blob

import (
	"crypto/sha256"
	"sync"

	"github.com/dustin/go-humanize"
)

// default blob store implementation satisfied Store interface

var chunkPool = sync.Pool{
	New: func() any {
		chunk := make([]byte, 16*humanize.MiByte)
		return &chunk
	},
}

var hashPool = sync.Pool{
	New: func() any {
		hash := sha256.New()
		return hash
	},
}

// // encapsulate the read and write file-system operation
// type localFS struct {
// 	root     string
// 	indexSUM map[string][]string
// }

// func newLocalFS(vol string) (*localFS, error) {
// 	if err := os.MkdirAll(vol, os.ModePerm); err != nil {
// 		return nil, err
// 	}

// 	l := &localFS{
// 		root:     vol,
// 		indexSUM: map[string][]string{},
// 	}
// 	return l, nil
// }

// func (lfs *localFS) Close() error {
// 	return nil
// }

// func (lfs *localFS) Put(ctx context.Context, objName string, size int64, r io.Reader) error {
// 	chunk := *chunkPool.Get().(*[]byte) //make([]byte, 16*humanize.MiByte)
// 	defer func() {
// 		clear(chunk)
// 		chunkPool.Put(&chunk)
// 	}()

// 	hash := hashPool.Get().(hash.Hash) //sha256.New()
// 	defer func() {
// 		hash.Reset()
// 		hashPool.Put(hash)
// 	}()

// 	for {
// 		n, err := r.Read(chunk)
// 		if err != nil {
// 			if err != io.EOF {
// 				return err
// 			}
// 			if n == 0 {
// 				//maybe eof
// 				return nil
// 			}

// 		}

// 		hash.Write(chunk[:n])
// 		key := hex.EncodeToString(hash.Sum(nil))

// 		path := filepath.Join(lfs.root, key)
// 		if err := os.WriteFile(path, chunk[:n], os.ModePerm); err != nil {
// 			return err
// 		}
// 		// pr.Size += int64(n)
// 		lfs.indexSUM[objName] = append(lfs.indexSUM[objName], key)
// 	}

// }

// func (lfs *localFS) Get(ctx context.Context, ri *ObjectInfo) (io.ReadCloser, uint32, error) {
// 	files := []*os.File{}
// 	var getErr error
// 	parts := lfs.indexSUM[ri.Filename]
// 	for _, part := range parts {
// 		path := filepath.Join(lfs.root, part)
// 		f, err := os.Open(path)
// 		if err != nil {
// 			getErr = err
// 			break
// 		}

// 		files = append(files, f)
// 	}

// 	if getErr != nil {
// 		return nil, 0, getErr
// 	}

// 	mr := mergeFileReader(files)
// 	fi := fileIterator{
// 		mr:    mr,
// 		files: files,
// 	}

// 	return &fi, 0, nil
// }

// func mergeFileReader(files []*os.File) io.Reader {
// 	readers := []io.Reader{}
// 	for _, f := range files {
// 		readers = append(readers, f)
// 	}

// 	return io.MultiReader(readers...)
// }

// type fileIterator struct {
// 	mr    io.Reader // io.multiReader
// 	files []*os.File
// }

// func (fi *fileIterator) Read(p []byte) (int, error) {
// 	return fi.mr.Read(p)
// }

// func (fi *fileIterator) Close() error {
// 	for _, f := range fi.files {
// 		err := f.Close()
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
