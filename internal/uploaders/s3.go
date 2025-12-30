package uploaders

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileUploader interface {
	RelativePath(path string) string
	RemoveFile(ctx context.Context, path string) error
	UploadFile(ctx context.Context, path string) error
}

type S3FileUploader struct {
	basePath	string
	bucket		string
	s3Client	*s3.Client
}

func NewS3FileUploader(s3Client *s3.Client, bucket string, basePath string) *S3FileUploader {
	return &S3FileUploader{
		basePath:	basePath,
		bucket:		bucket,
		s3Client:	s3Client,
	}
}

// UploadFile uploads a file from a given absolute path to the S3 bucket
// specified by the AWS_BUCKET env variable
func (u *S3FileUploader) UploadFile(ctx context.Context, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("os open: %v", err)
	}

	_, err = u.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(u.RelativePath(path)),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("s3 client putobject: %v", err)
	}

	return nil
}

// RemoveFile removes a file from a given absolute path from the S3 bucket
// specified by the AWS_BUCKET env variable
func (u *S3FileUploader) RemoveFile(ctx context.Context, path string) error {
	_, err := u.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(u.RelativePath(path)),
	})
	if err != nil {
		return fmt.Errorf("s3 client deleteobject: %v", err)
	}

	return nil
}

// RelativePath returns the relative path of a file from a given aboslute path string
func (u *S3FileUploader) RelativePath(path string) string {
	return strings.Trim(strings.Replace(path, u.basePath, "", 1), "/")
}
