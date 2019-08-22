package s3

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/viant/afs/storage"
	"os"
	"strings"
)

//Create creates a resource
func (s *storager) Create(ctx context.Context, destination string, mode os.FileMode, content []byte, isDir bool, options ...storage.Option) error {
	destination = strings.Trim(destination, "/")
	if !isDir {
		return s.Upload(ctx, destination, mode, content, options...)
	}
	return nil
}

func (s *storager) createBucket(ctx context.Context) error {
	_, err := s.S3.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: &s.bucket,
	})
	return err
}
