package rpc

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/odit-bit/cloudfs/internal/user"
	"github.com/odit-bit/cloudfs/internal/user/userpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

func init() {
	ctx, cancel := context.WithTimeout(context.TODO(), (2 * time.Second))
	_ = ctx
	defer cancel()
	//
	acc := user.NewWithMemory()
	svc := AuthService{
		accounts: acc,
		// logger:   logrus.New(),
	}
	lis = bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	userpb.RegisterAuthServiceServer(srv, &svc)
	go func() {
		if err := srv.Serve(lis); err != nil {
			if err.Error() == "closed" {
				return
			}
			log.Fatalf("error %v", err)
		}
	}()
}

// implement grpc dialer
func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func Test_auth(t *testing.T) {
	defer lis.Close()

	ctx := context.TODO()
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := userpb.NewAuthServiceClient(conn)

	username := "my-user-01"
	password := "my-password-01"
	testAccount, err := client.Register(ctx, &userpb.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatal(err)
	}

	resAuth, err := client.BasicAuth(ctx, &userpb.BasicAuthRequest{Username: username, Password: password})
	if err != nil {
		t.Fatal(err)
	}

	resTokenAuth, err := client.TokenAuth(ctx, &userpb.TokenAuthRequest{Token: resAuth.Token})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testAccount.UserID, resTokenAuth.UserID)
	assert.Equal(t, time.Now().Add(24*7*time.Hour).UTC().Round(24*time.Hour), resTokenAuth.ValidUntil.AsTime().Round(24*time.Hour))

}
