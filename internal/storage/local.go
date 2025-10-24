package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LocalStorage struct{ BasePath string; PublicBaseURL string }

func NewLocalStorage(basePath, publicBaseURL string) *LocalStorage { return &LocalStorage{BasePath: basePath, PublicBaseURL: publicBaseURL} }

func (s *LocalStorage) Save(ctx context.Context, r io.Reader, objectPath, contentType string) (string, error) {
	full := filepath.Join(s.BasePath, objectPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil { return "", err }
	f, err := os.Create(full)
	if err != nil { return "", err }
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil { return "", err }
	// best-effort set extension if missing
	if filepath.Ext(full) == "" && contentType != "" {
		if exts, _ := mime.ExtensionsByType(contentType); len(exts) > 0 {
			_ = os.Rename(full, full+exts[0])
			full = full + exts[0]
		}
	}
	return objectPath, nil
}

func (s *LocalStorage) SignedURL(ctx context.Context, objectPath string, expiry time.Duration) (string, error) {
	if s.PublicBaseURL != "" {
		return strings.TrimRight(s.PublicBaseURL, "/") + "/" + strings.TrimLeft(objectPath, "/"), nil
	}
	// fallback plain relative path
	return "/" + strings.TrimLeft(objectPath, "/"), nil
}

func (s *LocalStorage) Delete(ctx context.Context, objectPath string) error {
	full := filepath.Join(s.BasePath, objectPath)
	return os.Remove(full)
}

func (s *LocalStorage) BuildPath(parts ...string) string { return filepath.ToSlash(filepath.Join(parts...)) }
