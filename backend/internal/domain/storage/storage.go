package storage

import (
	"context"
	"io"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=./mock/storage_mock.go -package=mock
// Storage defines the interface for file storage operations
type Storage interface {
	// Upload uploads data to the storage and returns the object key/path
	Upload(ctx context.Context, key string, data io.Reader, contentType string) (string, error)

	// GetURL returns the accessible URL for the given object key
	GetURL(ctx context.Context, key string) (string, error)

	// GeneratePresignedURL generates a presigned URL for uploading
	GeneratePresignedURL(ctx context.Context, key string, contentType string, contentLength int64, expires time.Duration) (string, error)
}
