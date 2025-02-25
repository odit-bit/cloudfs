package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/blob/blobpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ blobpb.StorageServiceServer = (*BlobService)(nil)

type BlobService struct {
	objects *blob.Storage
	logger  *logrus.Logger
	blobpb.UnimplementedStorageServiceServer
}

func NewBlobService(objects *blob.Storage, logger *logrus.Logger) *BlobService {
	return &BlobService{
		objects: objects,
		logger:  logger,
	}
}

// DeleteObject implements blobpb.StorageServiceServer.
func (b *BlobService) DeleteObject(ctx context.Context, req *blobpb.DeleteRequest) (*blobpb.DeleteResponse, error) {

	if err := b.objects.Delete(ctx, req.Bucket, req.Filename); err != nil {
		return nil, err
	}

	return &blobpb.DeleteResponse{
		DeleteAt: timestamppb.New(time.Now().UTC()),
	}, nil
}

// DownloadObject implements blobpb.StorageServiceServer.
func (b *BlobService) DownloadObject(req *blobpb.DownloadRequest, stream grpc.ServerStreamingServer[blobpb.DownloadResponse]) error {

	info, err := b.objects.Download(stream.Context(), req.Bucket, req.Filename)
	if err != nil {
		return status.Error(codes.Aborted, err.Error())
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

		if err := stream.Send(&blobpb.DownloadResponse{
			TotalSize: info.Size,
			Chunk:     chunk[:n],
		}); err != nil {
			return err
		}

	}
	return nil
}

// DownloadSharedObject implements blobpb.StorageServiceServer.
func (b *BlobService) DownloadSharedObject(req *blobpb.DownloadSharedRequest, stream grpc.ServerStreamingServer[blobpb.DownloadResponse]) error {
	info, err := b.objects.DownloadToken(stream.Context(), req.SharedToken)
	if err != nil {
		if errors.Is(err, blob.ErrInvalidShareToken) {
			return status.Error(codes.Unauthenticated, err.Error())
		}
		if errors.Is(err, blob.ErrTokenExpired) {
			return status.Error(codes.Unauthenticated, err.Error())
		}

		return status.Error(codes.Internal, err.Error())
	}

	header := map[string]string{
		"filename":     info.Filename,
		"content-type": info.ContentType,
		"size":         strconv.FormatInt(info.Size, 10),
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

		if err := stream.Send(&blobpb.DownloadResponse{
			// TotalSize: info.Size,
			Chunk: chunk[:n],
		}); err != nil {
			return err
		}

	}

	return nil
}

// ShareObject implements blobpb.StorageServiceServer.
func (b *BlobService) ShareObject(ctx context.Context, req *blobpb.ShareObjectRequest) (*blobpb.ShareObjectResponse, error) {
	shareTkn, err := b.objects.CreateShareToken(ctx, req.Bucket, req.Filename, 0)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if shareTkn == nil {
		return nil, status.Error(codes.Unknown, "shareToken pointer is nil")
	}

	return &blobpb.ShareObjectResponse{
		ShareToken: shareTkn.Key,
		ValidUntil: timestamppb.New(shareTkn.ValidUntil()),
	}, nil
}

// UploadObject implements blobpb.StorageServiceServer.
func (b *BlobService) UploadObject(stream grpc.ClientStreamingServer[blobpb.UploadRequest, blobpb.UploadResponse]) error {

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Error(codes.Aborted, "ho header found")
	}

	var header uploadHeader
	if err := bindUploadHeader(md, &header); err != nil {
		return err
	}

	b.logger.Info("size:", header.totalSize)

	// uploading
	var written int
	pr, pw := io.Pipe()
	// cw := g.objects.UploadChunk(stream.Context(), tkn.UserID, header.filename, header.contentType, header.totalSize)
	go func() error {
		defer pw.Close()
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
			if nn, err := pw.Write(chunk); err != nil {
				return err
			} else {
				written += nn
			}
		}
		return nil
	}()

	res, err := b.objects.Put(stream.Context(),
		&blob.UploadParam{
			Bucket:      header.bucket,
			Filename:    header.filename,
			ContentType: header.contentType,
			Size:        header.totalSize,
			Body:        pr,
		})

	if err != nil {
		b.logger.Infof("isContextErr ? : %v", stream.Context().Err())
		b.logger.Errorf("failed put object: %v , written %v", err, written)
		return err
	}

	// if res.Size != header.totalSize {
	// 	return status.Error(codes.InvalidArgument, fmt.Sprintf("bytes written not match with TotalSize request: %d", written))
	// }

	if err := stream.SendAndClose(&blobpb.UploadResponse{
		Sum: res.Sum,
	}); err != nil {
		return err
	}

	b.logger.Infof("written bytes: %d", written)
	return nil
}

type uploadHeader struct {
	bucket      string
	filename    string
	contentType string
	totalSize   int64
}

func bindUploadHeader(md metadata.MD, header *uploadHeader) error {
	if len(md.Get("filename")) == 0 {
		return status.Error(codes.FailedPrecondition, "missing header filename")
	} else {
		header.filename = md.Get("filename")[0]
	}

	if len(md.Get("bucket")) == 0 {
		return status.Error(codes.FailedPrecondition, "missing header bucket")
	} else {
		header.bucket = md.Get("bucket")[0]
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

	return nil
}

/// Object Service

func (g *BlobService) ListObject(req *blobpb.ListObjectRequest, stream grpc.ServerStreamingServer[blobpb.ListObjectResponse]) error {
	if req.Bucket == "" {
		return status.Error(codes.InvalidArgument, "invalid bucket")
	}

	iter := g.objects.List(stream.Context(), req.Bucket, int(req.Limit), "")
	// if err != nil {
	// 	return status.Error(codes.Internal, err.Error())
	// }

	count := 0
	for iter.Next() {
		obj := iter.Value()
		count++
		if err := stream.Send(&blobpb.ListObjectResponse{
			Filename:     obj.Filename,
			Size:         obj.Size,
			ContentType:  obj.ContentType,
			Sum:          obj.Sum,
			LastModified: timestamppb.New(obj.LastModified),
		}); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	if iter.Err() != nil {
		return status.Error(codes.Internal, iter.Err().Error())
	}

	return nil

}
