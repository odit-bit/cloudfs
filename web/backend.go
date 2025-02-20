package web

import (
	"github.com/odit-bit/cloudfs/server/apipb"
	"google.golang.org/grpc"
)

// wrap backend grpc api call
type backend struct {
	cli apipb.StorageServiceClient
}

func NewBackend(conn *grpc.ClientConn) *backend {
	cli := apipb.NewStorageServiceClient(conn)
	return &backend{
		cli: cli,
	}
}

// func (b *backend) BasicAuth(ctx context.Context, username, password string) (*basicAuthResponse, error) {
// 	res, err := b.cli.BasicAuth(ctx, &apipb.BasicAuthRequest{
// 		Username: username,
// 		Password: password,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &basicAuthResponse{
// 		Token: res.Token,
// 	}, nil
// }

// func (b *backend) Register(ctx context.Context, username, password string) (*registerResponse, error) {

// }
