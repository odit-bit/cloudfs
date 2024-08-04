package s4

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/service"
)

var _ service.BlobStore = (*S3blobRepo)(nil)

type S3blobRepo struct {
	cli *s3.Client
}

// Delete implements service.BlobStore.
func (s *S3blobRepo) Delete(ctx context.Context, bucket string, filename string) error {
	res, err := s.cli.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &filename,
	})
	_ = res

	if err != nil {
		return err
	}
	return nil
}

// Get implements service.BlobStore.
func (s *S3blobRepo) Get(ctx context.Context, bucket string, filename string) (*blob.ObjectInfo, error) {
	panic("unimplemented")
}

// ObjectIterator implements service.BlobStore.
func (s *S3blobRepo) ObjectIterator(ctx context.Context, bucket string, limit int, lastFilename string) blob.Iterator {
	panic("unimplemented")
}

// Put implements service.BlobStore.
func (s *S3blobRepo) Put(ctx context.Context, bucket string, filename string, reader io.Reader, size int64, contentType string) (*blob.ObjectInfo, error) {
	panic("unimplemented")
}

func NewS3blobRepo(endpoint string, credential aws.Credentials) (*S3blobRepo, error) {
	// aws.Credentials{
	// 	AccessKeyID:     "access-key-id-123",
	// 	SecretAccessKey: "secret-access-key-123",
	// 	SessionToken:    "session-token-123",
	// 	Source:          "source-credential-langit",
	// 	CanExpire:       false,
	// 	Expires:         time.Time{},
	// 	AccountID:       "id-123-account",
	// }

	credsFunc := aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return credential, nil
	})

	svc := s3.New(s3.Options{
		BaseEndpoint:                   &endpoint,
		Credentials:                    credsFunc,
		DisableMultiRegionAccessPoints: true,
		DisableS3ExpressSessionAuth:    aws.Bool(true),
		Region:                         "langit",
		RetryMaxAttempts:               1,
		UsePathStyle:                   false,
	})

	return &S3blobRepo{
		cli: svc,
	}, nil
}
