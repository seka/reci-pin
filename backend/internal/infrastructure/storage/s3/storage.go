package s3

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
)

type client struct {
	client           *s3.Client
	presignClient    *s3.PresignClient
	bucket           string
	publicBaseURL    string
	internalEndpoint string // Internal endpoint for replacement logic
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
		client:           s3Client,
		presignClient:    s3.NewPresignClient(s3Client),
		bucket:           bucket,
		publicBaseURL:    publicBaseURL,
		internalEndpoint: endpoint,
	}, nil
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

	url := request.URL
	if c.publicBaseURL != "" && c.internalEndpoint != "" {
		// Replace internal endpoint + bucket with public base URL
		// Example: http://localstack:4566/recipin-bucket/ -> https://localhost/storage/
		internalBase := fmt.Sprintf("%s/%s/", c.internalEndpoint, c.bucket)
		publicBase := c.publicBaseURL
		if !strings.HasSuffix(publicBase, "/") {
			publicBase += "/"
		}
		url = strings.Replace(url, internalBase, publicBase, 1)
	}

	return url, nil
}

