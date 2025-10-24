package storage

import (
	"context"
	"io"
	"path"
	"time"
)

type Storage interface {
	Save(ctx context.Context, r io.Reader, objectPath, contentType string) (string, error)
	SignedURL(ctx context.Context, objectPath string, expiry time.Duration) (string, error)
	Delete(ctx context.Context, objectPath string) error
	BuildPath(parts ...string) string
}

func Join(parts ...string) string { return path.Join(parts...) }
