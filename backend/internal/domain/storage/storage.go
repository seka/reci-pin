package storage

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

// Client defines the interface for file storage operations
//
//go:generate mockgen -source=$GOFILE -destination=./mock/storage_mock.go -package=mock
type Client interface {
	// GetPublicURL returns a public URL for a file
	GetPublicURL() *url.URL
	// GeneratePresignedURL generates a presigned URL for uploading
	GeneratePresignedURL(ctx context.Context, key string, contentType model.RecipeImageType, contentLength int64, expires time.Duration) (string, error)
	// Upload uploads a file to the storage
	Upload(ctx context.Context, key string, body io.Reader, contentType model.RecipeImageType) error
	// Delete deletes a file from the storage
	Delete(ctx context.Context, key string) error
}
