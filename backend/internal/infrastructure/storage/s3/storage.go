package s3

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
)

type client struct {
	client             *s3.Client
	presignClient      *s3.PresignClient
	bucket             string
	publicBaseURL      *url.URL
	internalEndpoint   *url.URL
	internalPathPrefix string
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

	parsedPublic, err := url.Parse(publicBaseURL)
	if err != nil && publicBaseURL != "" {
		return nil, fmt.Errorf("invalid public base URL: %w", err)
	}

	parsedInternal, err := url.Parse(endpoint)
	if err != nil && endpoint != "" {
		return nil, fmt.Errorf("invalid internal endpoint: %w", err)
	}

	var internalPathPrefix string
	if parsedInternal != nil {
		internalPathPrefix = parsedInternal.JoinPath(bucket).Path
		if !strings.HasSuffix(internalPathPrefix, "/") {
			internalPathPrefix += "/"
		}
	}

	return &client{
		client:             s3Client,
		presignClient:      s3.NewPresignClient(s3Client),
		bucket:             bucket,
		publicBaseURL:      parsedPublic,
		internalEndpoint:   parsedInternal,
		internalPathPrefix: internalPathPrefix,
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

	urlStr := request.URL
	if c.publicBaseURL != nil && c.internalEndpoint != nil {
		parsedPresigned, err := url.Parse(urlStr)
		if err == nil {
			if after, ok := strings.CutPrefix(parsedPresigned.Path, c.internalPathPrefix); ok {
				relPath := after

				// Construct final public URL
				parsedPresigned.Scheme = c.publicBaseURL.Scheme
				parsedPresigned.Host = c.publicBaseURL.Host
				parsedPresigned.Path = c.publicBaseURL.JoinPath(relPath).Path

				urlStr = parsedPresigned.String()
			}
		}
	}

	return urlStr, nil
}
