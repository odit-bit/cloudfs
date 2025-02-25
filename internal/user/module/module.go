package module

import (
	"context"
	"database/sql"

	"github.com/odit-bit/cloudfs/internal/user"
	"github.com/odit-bit/cloudfs/internal/user/repo"
	"github.com/odit-bit/cloudfs/internal/user/rpc"
	"github.com/odit-bit/cloudfs/internal/user/userpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(ctx context.Context, logger *logrus.Logger, grpcSrv *grpc.Server, db *sql.DB) error {
	upg, err := repo.NewUserPG(ctx, db)
	if err != nil {
		return err
	}
	utpg, err := repo.NewUserTokenPG(ctx, db)
	if err != nil {
		return err
	}

	users, err := user.New(ctx, upg, utpg)
	if err != nil {
		return err
	}
	svc := rpc.NewAuthService(users, logger)
	userpb.RegisterAuthServiceServer(grpcSrv, svc)
	return nil
}
