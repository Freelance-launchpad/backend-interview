package js3

import (
	"context"
	"io"
)

type FileProvider interface {
	Store(ctx context.Context, identifier string, readSeeker io.ReadSeeker) error
	StoreWithMetadata(ctx context.Context, identifier string, metadata map[string]string, readSeeker io.ReadSeeker) error
	StoreWithMetadataAndContentType(ctx context.Context, identifier, contentType string, metadata map[string]string, readSeeker io.ReadSeeker) error
	Delete(ctx context.Context, identifier string) error
	GetURL(ctx context.Context, identifier string) (string, error)
	GetURLWithMetadata(ctx context.Context, identifier string) (string, map[string]string, error)
	Download(ctx context.Context, identifier string) (FileOutput, error)
}

type FileOutput struct {
	FileName      string
	Reader        io.ReadCloser
	ContentLength int64
	ContentType   string
}
