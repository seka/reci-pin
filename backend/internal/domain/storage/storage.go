package storage

import (
	"context"
	"net/url"
	"time"
)

// Client defines the interface for file storage operations
//
//go:generate mockgen -source=$GOFILE -destination=./mock/storage_mock.go -package=mock
type Client interface {
	// GetPublicURL returns a public URL for a file
	GetPublicURL() *url.URL
	// GeneratePresignedURL generates a presigned URL for uploading
	GeneratePresignedURL(ctx context.Context, key string, contentType string, contentLength int64, expires time.Duration) (string, error)
}
