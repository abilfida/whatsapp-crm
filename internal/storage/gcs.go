package storage

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

type GCSStorage struct{ Bucket, Prefix, GoogleAccessID string; PrivateKey []byte; Client *storage.Client }

func NewGCSStorage(ctx context.Context, bucket, prefix string, credentialsPath string) (*GCSStorage, error) {
	// Read service account file
	credsBytes, err := os.ReadFile(credentialsPath)
	if err != nil { return nil, err }
	cfg, err := google.JWTConfigFromJSON(credsBytes)
	if err != nil { return nil, err }
	// Extract fields for signing
	accessID := cfg.Email
	// Parse private key (PKCS#8 or PKCS#1)
	var keyBytes []byte
	if block, _ := pem.Decode(credsBytes); block != nil {
		keyBytes = block.Bytes
	} else {
		keyBytes = cfg.PrivateKey
	}

	client, err := storage.NewClient(ctx)
	if err != nil { return nil, err }
	return &GCSStorage{Bucket: bucket, Prefix: prefix, GoogleAccessID: accessID, PrivateKey: keyBytes, Client: client}, nil
}

func (g *GCSStorage) obj(p string) *storage.ObjectHandle { key := p; if g.Prefix != "" { key = g.Prefix + "/" + p }; return g.Client.Bucket(g.Bucket).Object(key) }

func (g *GCSStorage) Save(ctx context.Context, r io.Reader, objectPath, contentType string) (string, error) {
	w := g.obj(objectPath).NewWriter(ctx)
	w.ContentType = contentType
	if _, err := io.Copy(w, r); err != nil { _ = w.Close(); return "", err }
	if err := w.Close(); err != nil { return "", err }
	return objectPath, nil
}

func (g *GCSStorage) SignedURL(ctx context.Context, objectPath string, expiry time.Duration) (string, error) {
	url, err := storage.SignedURL(g.Bucket, g.key(objectPath), &storage.SignedURLOptions{
		GoogleAccessID: g.GoogleAccessID,
		PrivateKey:     g.PrivateKey,
		Method:         "GET",
		Expires:        time.Now().Add(expiry),
	})
	return url, err
}

func (g *GCSStorage) key(p string) string { if g.Prefix == "" { return p }; return g.Prefix + "/" + p }

func (g *GCSStorage) Delete(ctx context.Context, objectPath string) error { return g.obj(objectPath).Delete(ctx) }

func (g *GCSStorage) BuildPath(parts ...string) string { return Join(parts...) }
