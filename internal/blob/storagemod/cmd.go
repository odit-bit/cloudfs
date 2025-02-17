package storagemodule

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/blob/storagegrpc"
	"github.com/odit-bit/cloudfs/internal/blob/storagepb"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, mux *chi.Mux, rpc *grpc.Server) error {
	storage, err := blob.NewWithMemory()
	if err != nil {
		return err
	}
	storagepb.RegisterStorageServiceServer(rpc, storagegrpc.New(storage))
	return nil
}
