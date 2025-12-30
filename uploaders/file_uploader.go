package uploaders

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// UploadFile uploads a file from a given absolute path to the S3 bucket
// specified by the AWS_BUCKET env variable
func UploadFile(path string) (*s3.PutObjectOutput, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("could not open path")
	}

	return s3Client().PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(RelativePath(path)),
		Body:   file,
	})
}

// RemoveFile removes a file from a given absolute path from the S3 bucket
// specified by the AWS_BUCKET env variable
func RemoveFile(path string) (*s3.DeleteObjectOutput, error) {
	return s3Client().DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(RelativePath(path)),
	})
}

// RelativePath returns the relative path of a file from a given aboslute path string
func RelativePath(path string) string {
	return strings.Trim(strings.Replace(path, os.Getenv("SINKER_BASE_PATH"), "", 1), "/")
}

// s3Client returns a new S3 client
func s3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = "eu-west-1"
	})
}
