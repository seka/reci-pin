package s3

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	sdkconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/seka/reci-pin/backend/config"
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
func NewClient(ctx context.Context, cfg config.Storage) (storage.Client, error) {
	awsCfg, err := sdkconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true // Required for LocalStack
		}
	})

	var parsedPublic *url.URL
	if cfg.PublicBaseURL != "" {
		parsedPublic, err = url.Parse(cfg.PublicBaseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid public base URL: %w", err)
		}
	}

	var parsedInternal *url.URL
	if cfg.Endpoint != "" {
		parsedInternal, err = url.Parse(cfg.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid internal endpoint: %w", err)
		}
	}

	var internalPathPrefix string
	if parsedInternal != nil {
		internalPathPrefix = parsedInternal.JoinPath(cfg.Bucket).Path
		if !strings.HasPrefix(internalPathPrefix, "/") {
			internalPathPrefix = "/" + internalPathPrefix
		}
		if !strings.HasSuffix(internalPathPrefix, "/") {
			internalPathPrefix += "/"
		}
	}

	return &client{
		client:             s3Client,
		presignClient:      s3.NewPresignClient(s3Client),
		bucket:             cfg.Bucket,
		publicBaseURL:      parsedPublic,
		internalEndpoint:   parsedInternal,
		internalPathPrefix: internalPathPrefix,
	}, nil
}

func (c *client) GetPublicURL() *url.URL {
	return c.publicBaseURL
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
