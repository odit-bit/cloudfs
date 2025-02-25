package module

import (
	"context"
	"database/sql"

	"github.com/minio/minio-go/v7"
	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/blob/blobpb"
	"github.com/odit-bit/cloudfs/internal/blob/repo"
	"github.com/odit-bit/cloudfs/internal/blob/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(ctx context.Context, logger *logrus.Logger, cli *minio.Client, sql *sql.DB, grpcSrv *grpc.Server) error {
	objects := repo.NewMinioBlob(cli)
	objectToken, err := repo.NewPGShareToken(ctx, sql)
	if err != nil {
		return err
	}

	storage, err := blob.New(ctx, objectToken, objects)
	if err != nil {
		return err
	}

	svc := rpc.NewBlobService(storage, logger)
	blobpb.RegisterStorageServiceServer(grpcSrv, svc)

	return nil
}
