package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/odit-bit/cloudfs/rpc/apipb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// wrapper for grpc client

// wrap backend grpc api call
type CloudfsClient struct {
	cli apipb.StorageServiceClient
}

func NewCloudfsClient(conn *grpc.ClientConn) *CloudfsClient {
	cli := apipb.NewStorageServiceClient(conn)
	return &CloudfsClient{
		cli: cli,
	}
}

////

type BasicAuthResponse struct {
	Token string
}

func (b *CloudfsClient) BasicAuth(ctx context.Context, username, password string) (*BasicAuthResponse, error) {
	res, err := b.cli.BasicAuth(ctx, &apipb.BasicAuthRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &BasicAuthResponse{
		Token: res.Token,
	}, nil
}

////

type RegisterParam struct {
	Username string
	Password string
}

type RegisterResult struct {
	UserID string
}

func (b *CloudfsClient) Register(ctx context.Context, param RegisterParam) (*RegisterResult, error) {
	resp, err := b.cli.Register(ctx, &apipb.RegisterRequest{Username: param.Username, Password: param.Password})
	if err != nil {
		return nil, err
	}
	return &RegisterResult{UserID: resp.UserID}, nil
}

////

type UploadResult struct {
	Sum string
}

func (b *CloudfsClient) UploadObject(ctx context.Context, userToken, filename, contentType string, size int64, body io.Reader) (*UploadResult, error) {
	md := metadata.New(map[string]string{})
	md.Set("filename", filename)
	md.Set("authorization", userToken)
	md.Set("content-type", contentType)
	md.Set("content-length", strconv.Itoa(int(size)))

	sCtx := metadata.NewOutgoingContext(ctx, md)
	stream, err := b.cli.UploadObject(sCtx)
	if err != nil {
		return nil, fmt.Errorf("send metadata:%v", err)
	}

	chunk := make([]byte, 1024*1024*3)
	written := 0
	for {
		// if n is not 0 try to process the bytes first,
		// process the io.EOF in the next iteration

		n, err := body.Read(chunk)
		if err != nil && n == 0 {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("read chunk: %v", err)
		}

		if n == 0 {
			return nil, fmt.Errorf("read chunk: readed byte is zero and but err is nil")
		}

		if err := stream.Send(&apipb.UploadRequest{
			Chunk: chunk[:n]},
		); err != nil {
			return nil, fmt.Errorf("upload chunk: %v", err)
		}
		written += n
	}
	log.Printf("grpc client written bytes: %d", written)
	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("close upload stream: %v", err)
	}

	return &UploadResult{
		Sum: res.Sum,
	}, nil
}

type GetObjectResult struct {
	Filename    string
	ContentType string
	Size        int64
	Reader      io.ReadCloser
}

func (b *CloudfsClient) DownloadWithToken(ctx context.Context, shareToken string) (*GetObjectResult, error) {
	stream, err := b.cli.DownloadSharedObject(ctx, &apipb.DownloadSharedRequest{
		SharedToken: shareToken,
	})
	if err != nil {
		return nil, err
	}

	md, err := stream.Header()
	if err != nil {
		return nil, fmt.Errorf("failed receive header: %v", err)
	}
	if md == nil {
		_, xerr := stream.Recv()
		err = errors.Join(err, xerr)
		return nil, err
	}

	var filename, contentType string
	if xfilename := md.Get("filename"); len(xfilename) == 0 {
		return nil, fmt.Errorf("missing header filename from server")
	} else {
		filename = xfilename[0]
	}
	if xct := md.Get("filename"); len(xct) == 0 {
		return nil, fmt.Errorf("missing header filename from server")
	} else {
		contentType = xct[0]
	}

	pr, pw := io.Pipe()
	res := &GetObjectResult{
		Filename:    filename,
		ContentType: contentType,
		Reader:      pr,
	}

	go func() {
		defer pw.Close()
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				stream.CloseSend()
				return
			}
			if _, err := pw.Write(res.Chunk); err != nil {
				return
			}
		}
	}()

	return res, nil
}

func (b *CloudfsClient) DownloadObject(ctx context.Context, userToken, filename string) (*GetObjectResult, error) {
	objStream, err := b.cli.DownloadObject(ctx, &apipb.DownloadRequest{
		Token:    userToken,
		Filename: filename,
	})
	if err != nil {
		return nil, err
	}

	pr, pw := io.Pipe()
	res := &GetObjectResult{
		Reader: pr,
	}

	go func() {
		defer pw.Close()
		for {
			res, err := objStream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				objStream.CloseSend()
				return
			}
			if _, err := pw.Write(res.Chunk); err != nil {
				objStream.CloseSend()
				return
			}
		}
	}()

	return res, nil
}

type Object struct {
	UserID       string
	Filename     string
	Sum          string
	Size         int64
	LastModified time.Time
	ContentType  string
	err          error
}

func (o *Object) Err() error {
	return o.err
}

func (b *CloudfsClient) Objects(ctx context.Context, userToken, lastFilename string, limit int) (<-chan *Object, error) {
	stream, err := b.cli.ListObject(ctx, &apipb.ListObjectRequest{
		UserToken:    userToken,
		Limit:        1000,
		LastFilename: lastFilename,
	})
	if err != nil {
		return nil, err
	}

	objC := make(chan *Object, 1)

	go func() {
		for {
			obj := &Object{}
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					obj.err = err
					stream.CloseSend()
				}

			} else {
				obj.UserID = res.UserID
				obj.Filename = res.Filename
				obj.Sum = res.Filename
				obj.Size = res.Size
				obj.LastModified = res.LastModified.AsTime()
			}

			select {
			case <-ctx.Done():
				return
			default:
				objC <- obj
			}
		}
		close(objC)
	}()

	return objC, nil

}

////

type SharedObject struct {
	ShareToken string
	ValidUntil time.Time
}

func (b *CloudfsClient) ShareObject(ctx context.Context, userToken, filename string) (*SharedObject, error) {
	res, err := b.cli.ShareObject(ctx, &apipb.ShareObjectRequest{
		Token:    userToken,
		Filename: filename,
	})
	if err != nil {
		return nil, err
	}

	return &SharedObject{ShareToken: res.ShareToken, ValidUntil: res.ValidUntil.AsTime()}, nil
}

////

type DeleteResponse struct {
	DeletedAt time.Time
}

func (b *CloudfsClient) Delete(ctx context.Context, userToken, filename string) (*DeleteResponse, error) {
	res, err := b.cli.DeleteObject(ctx, &apipb.DeleteRequest{
		UserToken: userToken,
		Filename:  filename,
	})
	if err != nil {
		return nil, err
	}
	return &DeleteResponse{
		DeletedAt: res.DeleteAt.AsTime(),
	}, nil
}
