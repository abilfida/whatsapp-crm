package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct{ Bucket, Prefix string; Client *s3.Client }

func NewS3Storage(region, bucket, prefix, accessKey, secret string) (*S3Storage, error) {
	var cfg aws.Config
	var err error
	if accessKey != "" && secret != "" {
		creds := aws.NewCredentialsCache(aws.StaticCredentialsProvider{Value: aws.Credentials{AccessKeyID: accessKey, SecretAccessKey: secret}})
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRegion(region), config.WithCredentialsProvider(creds))
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	}
	if err != nil { return nil, err }
	return &S3Storage{Bucket: bucket, Prefix: prefix, Client: s3.NewFromConfig(cfg)}, nil
}

func (s *S3Storage) key(p string) string { if s.Prefix == "" { return p }; return s.Prefix + "/" + p }

func (s *S3Storage) Save(ctx context.Context, r io.Reader, objectPath, contentType string) (string, error) {
	uploader := manager.NewUploader(s.Client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      &s.Bucket,
		Key:         aws.String(s.key(objectPath)),
		ACL:         types.ObjectCannedACLPrivate,
		Body:        r,
		ContentType: aws.String(contentType),
	})
	if err != nil { return "", err }
	return objectPath, nil
}

func (s *S3Storage) SignedURL(ctx context.Context, objectPath string, expiry time.Duration) (string, error) {
	pres := s3.NewPresignClient(s.Client)
	req, err := pres.PresignGetObject(ctx, &s3.GetObjectInput{Bucket: &s.Bucket, Key: aws.String(s.key(objectPath))}, s3.WithPresignExpires(expiry))
	if err != nil { return "", err }
	return req.URL, nil
}

func (s *S3Storage) Delete(ctx context.Context, objectPath string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &s.Bucket, Key: aws.String(s.key(objectPath))})
	return err
}

func (s *S3Storage) BuildPath(parts ...string) string { return Join(parts...) }
