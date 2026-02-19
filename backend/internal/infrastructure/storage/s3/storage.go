package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
)

type client struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	uploader      *manager.Uploader
	bucket        string
	publicBaseURL string
}

// NewClient creates a new StorageService backed by S3
func NewClient(ctx context.Context, bucket string, endpoint string, publicBaseURL string) (storage.Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true // Required for LocalStack
		}
	})

	return &client{
		client:        s3Client,
		presignClient: s3.NewPresignClient(s3Client),
		uploader:      manager.NewUploader(s3Client),
		bucket:        bucket,
		publicBaseURL: publicBaseURL,
	}, nil
}

func (c *client) Upload(ctx context.Context, key string, data io.Reader, contentType string) (string, error) {
	result, err := c.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        data,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	if c.publicBaseURL != "" {
		return fmt.Sprintf("%s/%s", c.publicBaseURL, key), nil
	}
	return result.Location, nil
}

func (c *client) GetURL(ctx context.Context, key string) (string, error) {
	// For now, we return the direct URL. In the future, we might want to generate presigned URLs if the bucket is private.
	// However, for public read buckets or LocalStack, constructing the URL might be needed manually if Location is not enough.
	// The Upload method returns Location, which is usually the URL.
	// But let's assume we want to construct it or return what's stored.

	// If the bucket is public, we can construct the URL.
	// For LocalStack, it might be http://localhost:4566/<bucket>/<key>
	// But since the frontend accesses it, it needs to be accessible from the browser.
	// LocalStack endpoint is localhost:4566, so it should be fine.

	// For now, let's just return the key as the identifier, and let the frontend construct the URL or have a separate valid implementation.
	// Actually, `Upload` returns `Location`.

	return key, nil
}

func (c *client) GeneratePresignedURL(ctx context.Context, key string, contentType string, contentLength int64, expires time.Duration) (string, error) {
	request, err := c.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
	}, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return request.URL, nil
}
