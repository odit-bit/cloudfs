package rpc

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/blob/blobpb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

func init() {
	ctx, cancel := context.WithTimeout(context.TODO(), (2 * time.Second))
	_ = ctx
	defer cancel()

	//
	s, _ := blob.NewWithMemory()
	svc := BlobService{
		objects: s,
		logger:  logrus.New(),
	}
	lis = bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	blobpb.RegisterStorageServiceServer(srv, &svc)
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

func Test_server(t *testing.T) {
	defer lis.Close()

	ctx := context.TODO()
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := blobpb.NewStorageServiceClient(conn)

	// UPLOAD
	data := []byte("this is data content")
	headerObj := map[string]string{
		"filename":       "my-file.txt",
		"bucket":         "my-bucket",
		"content-length": strconv.Itoa(len(data)),
		"content-type":   "text",
	}
	assertUploadObject(ctx, t, client, data, headerObj)

	//DOWNLOAD
	assertDownloadObject(ctx, t, client, data, headerObj)

	// CreateShareToken
	shareToken, err := client.ShareObject(ctx, &blobpb.ShareObjectRequest{
		Bucket:   headerObj["bucket"],
		Filename: headerObj["filename"],
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = shareToken

	// Download With Share Token
	stream, err := client.DownloadSharedObject(ctx, &blobpb.DownloadSharedRequest{
		SharedToken: shareToken.ShareToken,
	})
	if err != nil {
		t.Fatal(err)
	}

	data2 := []byte{}
	for {
		v, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
			break
		}
		data2 = append(data2, v.GetChunk()...)
	}
	stream.CloseSend()
	assert.Equal(t, data, data2)

	//DELETE
	now := time.Now()
	req, err := client.DeleteObject(ctx, &blobpb.DeleteRequest{
		Bucket:   headerObj["bucket"],
		Filename: headerObj["filename"],
	})
	if err != nil {
		t.Fatal(err)
	}
	if req.DeleteAt.AsTime().Before(now) {
		t.Fatal()
	}
}

func assertDownloadObject(ctx context.Context, t *testing.T, client blobpb.StorageServiceClient, data []byte, headerObj map[string]string) {
	dlMD := metadata.New(map[string]string{})
	dlCtx := metadata.NewOutgoingContext(ctx, dlMD)
	stream2, err := client.DownloadObject(dlCtx, &blobpb.DownloadRequest{
		Bucket:   headerObj["bucket"],
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
}

func assertUploadObject(ctx context.Context, t *testing.T, client blobpb.StorageServiceClient, data []byte, headerObj map[string]string) {
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
		if err := stream.Send(&blobpb.UploadRequest{Chunk: data[start:end]}); err != nil {
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

}
