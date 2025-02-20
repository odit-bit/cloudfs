package server

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/odit-bit/cloudfs/server/apipb"
	"github.com/odit-bit/cloudfs/internal/storage"
	"github.com/odit-bit/cloudfs/internal/user"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

var account *user.Account
var token string

func init() {
	//
	ctx, cancel := context.WithTimeout(context.TODO(), (2 * time.Second))
	_ = ctx
	defer cancel()

	h := mockServer()
	//register user
	acc, err := h.accounts.Register(ctx, "mock-user", "mock-password")
	if err != nil {
		log.Fatal(err)
	}
	account = acc

	// auth user
	tkn, err := h.accounts.BasicAuth(ctx, "mock-user", "mock-password")
	if err != nil {
		log.Fatal(err)
	}
	token = tkn.Token.Key()

	//
	lis = bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	apipb.RegisterStorageServiceServer(srv, h)
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

func mockServer() *GrpcServer {
	s, _ := storage.NewWithMemory()
	u := user.NewWithMemory()

	srv := GrpcServer{
		objects:  s,
		accounts: u,
	}
	return &srv
}

func Test_Object(t *testing.T) {
	defer lis.Close()

	ctx := context.TODO()
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := apipb.NewStorageServiceClient(conn)

	username := "my-user-01"
	password := "my-password-01"
	testAccount, err := client.Register(ctx, &apipb.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatal(err)
	}

	resAuth, err := client.BasicAuth(ctx, &apipb.BasicAuthRequest{Username: username, Password: password})
	if err != nil {
		t.Fatal(err)
	}

	resTokenAuth, err := client.TokenAuth(ctx, &apipb.TokenAuthRequest{Token: resAuth.Token})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testAccount.UserID, resTokenAuth.UserID)
	assert.Equal(t, time.Now().Add(24*7*time.Hour).UTC().Round(24*time.Hour), resTokenAuth.ValidUntil.AsTime().Round(24*time.Hour))

	//
	//
	//
	//
	//

	// UPLOAD
	data := []byte("this is data content")
	headerObj := map[string]string{
		"filename":       "my-file.txt",
		"authorization":  resAuth.Token,
		"content-length": strconv.Itoa(len(data)),
		"content-type":   "text",
	}
	sCtx := metadata.NewOutgoingContext(ctx, metadata.New(headerObj))
	stream, err := client.UploadObject(sCtx)
	if err != nil {
		if err.Error() != "closed" {
			t.Fatal(err)
		}
	}

	start := 0
	end := 1
	for {
		if end > len(data) {
			break
		}

		// obj.Chunk =
		if err := stream.Send(&apipb.UploadRequest{Chunk: data[start:end]}); err != nil {
			t.Fatal(err)
		}
		start = end
		end += 1
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}
	if res.Sum == "" {
		t.Fatal("sum cannot empty")
	}

	//DOWNLOAD
	dlMD := metadata.New(map[string]string{"authorization": headerObj["authorization"]})
	dlCtx := metadata.NewOutgoingContext(ctx, dlMD)
	stream2, err := client.DownloadObject(dlCtx, &apipb.DownloadRequest{
		Token:    headerObj["authorization"],
		Filename: headerObj["filename"],
	})
	if err != nil {
		t.Fatal(err)
	}

	if in, err := stream2.Header(); err != nil {
		t.Fatal(err)
	} else {
		if in == nil {
			t.Fatal("incoming header from server is nil")
		}
		assert.Equal(t, in.Get("content-length")[0], headerObj["content-length"])
	}

	data2 := []byte{}
	isLoop := false
	for {
		isLoop = true
		resp, err := stream2.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		b := resp.GetChunk()
		data2 = append(data2, b...)
	}
	assert.Equal(t, true, isLoop)
	assert.Equal(t, string(data), string(data2))

	if err := stream2.CloseSend(); err != nil {
		t.Fatal(err)
	}
	//
	///
	//
	//
	//
	//
	//
	//

	resShareTkn, err := client.ShareObject(ctx, &apipb.ShareObjectRequest{
		Token:    headerObj["authorization"],
		Filename: headerObj["filename"],
	})
	if err != nil {
		t.Fatal(err)
	}

	stream3, err := client.DownloadSharedObject(ctx, &apipb.DownloadSharedRequest{
		SharedToken: resShareTkn.ShareToken,
	})
	if err != nil {
		t.Fatal(err)
	}

	data3 := []byte{}
	for {
		v, err := stream3.Recv()
		if err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
			break
		}
		data3 = append(data3, v.GetChunk()...)
	}
	stream3.CloseSend()
	assert.Equal(t, data2, data3)
}
