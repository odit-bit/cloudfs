package storagegrpc

import (
	"errors"
	"fmt"
	"io"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/blob/storagepb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var _ storagepb.StorageServiceServer = (*handler)(nil)

type handler struct {
	app *blob.Blobs
	storagepb.UnimplementedStorageServiceServer
}

func New(app *blob.Blobs) *handler {
	return &handler{
		app:                               app,
		UnimplementedStorageServiceServer: storagepb.UnimplementedStorageServiceServer{},
	}
}

// PutObject implements storagepb.StorageServiceServer.
func (s *handler) PutObject(stream grpc.ClientStreamingServer[storagepb.PutObjectRequest, storagepb.PutObjectResponse]) error {

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	if req.Info.Size <= 0 {
		return fmt.Errorf("cannot upload zero byte")
	}

	pr, pw := io.Pipe()
	eg := errgroup.Group{}
	eg.Go(func() error {
		defer pw.Close()
		if _, err := pw.Write(req.GetChunk().Chunk); err != nil {
			return err
		}

		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			if _, err := pw.Write(req.GetChunk().Chunk); err != nil {
				return err
			}
		}
	})

	_, uploadErr := s.app.Upload(stream.Context(), blob.UploadParam{
		Bucket:      req.Info.Bucket,
		Filename:    req.Info.Filename,
		ContentType: req.Info.ContentType,
		Size:        req.Info.Size,
		Body:        pr,
	})
	if err != nil {
		err = errors.Join(err, uploadErr)
	}

	wErr := eg.Wait()
	if wErr != nil {
		err = errors.Join(err, wErr)
	}

	err = errors.Join(err, stream.SendAndClose(&storagepb.PutObjectResponse{}))
	return err
}
