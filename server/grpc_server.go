package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/odit-bit/cloudfs/internal/storage"
	"github.com/odit-bit/cloudfs/internal/user"
	"github.com/odit-bit/cloudfs/server/apipb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ apipb.StorageServiceServer = (*GrpcServer)(nil)

type GrpcServer struct {
	objects  *storage.Blobs
	accounts *user.Users
	apipb.UnimplementedStorageServiceServer
}

func NewGrpcServer(objects *storage.Blobs, accounts *user.Users) *GrpcServer {
	return &GrpcServer{
		objects:                           objects,
		accounts:                          accounts,
		UnimplementedStorageServiceServer: apipb.UnimplementedStorageServiceServer{},
	}
}

func (g *GrpcServer) Register(ctx context.Context, req *apipb.RegisterRequest) (*apipb.RegisterResponse, error) {
	acc, err := g.accounts.Register(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &apipb.RegisterResponse{
		UserID: acc.ID.String(),
	}, nil
}

func (g *GrpcServer) BasicAuth(ctx context.Context, req *apipb.BasicAuthRequest) (*apipb.BasicAuthResponse, error) {
	res, err := g.accounts.BasicAuth(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &apipb.BasicAuthResponse{
		Token: res.Token.Key(),
	}, nil
}

func (g *GrpcServer) TokenAuth(ctx context.Context, req *apipb.TokenAuthRequest) (*apipb.TokenAuthResponse, error) {
	res, err := g.accounts.TokenAuth(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	return &apipb.TokenAuthResponse{
		UserID:     res.UserID(),
		ValidUntil: timestamppb.New(res.ValidUntil()),
	}, nil
}

/// Object Service

func (g *GrpcServer) ListObject(req *apipb.ListObjectRequest, stream grpc.ServerStreamingServer[apipb.ListObjectResponse]) error {
	if req.UserToken == "" {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	tkn, err := g.accounts.TokenAuth(stream.Context(), req.UserToken)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	list, err := g.objects.ListObject(stream.Context(), tkn.UserID(), int(req.Limit), "")
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, obj := range list {
		if err := stream.Send(&apipb.ListObjectResponse{
			Filename:     obj.Filename,
			Size:         obj.Size,
			ContentType:  obj.ContentType,
			Sum:          obj.Sum,
			LastModified: timestamppb.New(obj.LastModified),
		}); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil

}

func (g *GrpcServer) DeleteObject(ctx context.Context, req *apipb.DeleteRequest) (*apipb.DeleteResponse, error) {
	tkn, err := g.accounts.TokenAuth(ctx, req.UserToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if err := g.objects.Delete(ctx, tkn.UserID(), req.Filename); err != nil {
		return nil, err
	}

	return &apipb.DeleteResponse{
		DeleteAt: timestamppb.New(time.Now().UTC()),
	}, nil
}

func (g *GrpcServer) ShareObject(ctx context.Context, req *apipb.ShareObjectRequest) (*apipb.ShareObjectResponse, error) {
	tkn, err := g.accounts.TokenAuth(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	shareTkn, err := g.objects.CreateShareToken(ctx, tkn.UserID(), req.Filename, 24*time.Hour)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apipb.ShareObjectResponse{
		ShareToken: shareTkn.Key(),
		ValidUntil: timestamppb.New(shareTkn.ValidUntil()),
	}, nil
}

func (g *GrpcServer) DownloadSharedObject(req *apipb.DownloadSharedRequest, stream grpc.ServerStreamingServer[apipb.DownloadResponse]) error {
	info, err := g.objects.DownloadToken(stream.Context(), req.SharedToken)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidShareToken) {
			return status.Error(codes.Unauthenticated, err.Error())
		}
		if errors.Is(err, storage.ErrTokenExpired) {
			return status.Error(codes.Unauthenticated, err.Error())
		}
		if errors.Is(err, storage.ErrUnknown) {
			return status.Error(codes.Internal, err.Error())
		}
	}

	header := map[string]string{
		"filename":     info.Filename,
		"content-type": info.ContentType,
	}

	if err := stream.SendHeader(metadata.New(header)); err != nil {
		return status.Error(codes.Aborted, "failed send header")
	}

	defer info.Data.Close()
	chunk := make([]byte, 1024*1024*3)
	for {
		n, err := info.Data.Read(chunk)
		if err != nil {
			if err != io.EOF {
				return err
			}
		}

		if n == 0 {
			break
		}

		if err := stream.Send(&apipb.DownloadResponse{
			TotalSize: info.Size,
			Chunk:     chunk[:n],
		}); err != nil {
			return err
		}

	}

	return nil
}

// Download implements apipb.StorageServiceServer.
func (g *GrpcServer) DownloadObject(req *apipb.DownloadRequest, stream grpc.ServerStreamingServer[apipb.DownloadResponse]) error {
	// inMD, ok := metadata.FromIncomingContext(stream.Context())
	// if !ok {
	// 	return status.Error(codes.Aborted, "no header found")
	// }

	// var userToken string
	// if xAuth := inMD.Get("authorization"); len(xAuth) == 0 {
	// 	return status.Error(codes.Unauthenticated, "missing authorization header")
	// } else {
	// 	userToken = xAuth[0]
	// }

	tkn, err := g.accounts.TokenAuth(stream.Context(), req.Token)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}
	info, err := g.objects.Download(stream.Context(), tkn.UserID(), req.Filename)
	if err != nil {
		if errors.Is(err, storage.ErrUnknown) {
			return status.Error(codes.Aborted, err.Error())
		}
	}

	//send header before send first message
	outMD := metadata.New(map[string]string{"content-length": strconv.FormatInt(info.Size, 10)})
	if err := stream.SendHeader(outMD); err != nil {
		return status.Error(codes.Aborted, fmt.Sprintf("stream failed to send header response %v", err))
	}

	defer info.Data.Close()
	chunk := make([]byte, 1024*1024*3)
	for {
		n, err := info.Data.Read(chunk)
		if err != nil {
			if err != io.EOF {
				return err
			}
		}

		if n == 0 {
			break
		}

		if err := stream.Send(&apipb.DownloadResponse{
			TotalSize: info.Size,
			Chunk:     chunk[:n],
		}); err != nil {
			return err
		}

	}

	return nil

}

// Upload implements apipb.StorageServiceServer.
func (g *GrpcServer) UploadObject(stream grpc.ClientStreamingServer[apipb.UploadRequest, apipb.UploadResponse]) error {

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Error(codes.Aborted, "ho header found")
	}

	header := struct {
		filename    string
		token       string
		contentType string
		totalSize   int64
	}{}

	if len(md.Get("authorization")) == 0 {
		return status.Error(codes.FailedPrecondition, "missing header authorization")
	} else {
		header.token = md.Get("authorization")[0]
	}

	if len(md.Get("filename")) == 0 {
		return status.Error(codes.FailedPrecondition, "missing header filename")
	} else {
		header.filename = md.Get("filename")[0]
	}

	if ts := md.Get("content-length"); len(ts) == 0 {
		return status.Error(codes.FailedPrecondition, "missing header content-length")
	} else {
		totalSize, err := strconv.ParseInt(md.Get("content-length")[0], 10, 64)
		if err != nil {
			return status.Error(codes.Aborted, "invalid content-length value")
		}
		if totalSize <= 0 {
			return status.Error(codes.InvalidArgument, "size cannot be/below 0")
		}
		header.totalSize = totalSize
	}

	if ct := md.Get("content-type"); len(ct) != 0 {
		header.contentType = ct[0]
	}

	//authorize
	tkn, err := g.accounts.TokenAuth(stream.Context(), header.token)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	// uploading
	var written int
	cw := g.objects.NewChunkWriter(stream.Context(), tkn.UserID(), header.filename, header.contentType, header.totalSize)

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		chunk := req.GetChunk()
		if len(chunk) <= 0 {
			return status.Error(codes.Aborted, fmt.Sprintf("cannot write zero byte, err: %v", err))
		}

		if nn, err := cw.Write(chunk); err != nil {
			return err
		} else {
			written += nn
		}
	}

	res, err := cw.Result()
	if err != nil {
		return err
	}
	if res.Size != header.totalSize {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("bytes written not match with TotalSize request: %d", written))
	}

	if err := stream.SendAndClose(&apipb.UploadResponse{
		Sum: res.Sum,
	}); err != nil {
		return err
	}

	return nil
}
