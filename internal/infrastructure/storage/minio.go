package storage

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/agnathor/finances-go/internal/config"
	"github.com/agnathor/finances-go/internal/domain"
)

type minioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioClient(cfg config.MinIOConfig) (*minio.Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return client, nil
}

func NewStorageService(client *minio.Client, bucket string) domain.StorageService {
	return &minioStorage{
		client: client,
		bucket: bucket,
	}
}

func (s *minioStorage) Upload(ctx context.Context, filename, contentType string, data []byte) (string, error) {
	ext := filepath.Ext(filename)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	_, err := s.client.PutObject(ctx, s.bucket, uniqueName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload to minio: %w", err)
	}

	url := fmt.Sprintf("/api/v1/files/%s", uniqueName)

	return url, nil
}

func (s *minioStorage) Delete(ctx context.Context, filename string) error {
	err := s.client.RemoveObject(ctx, s.bucket, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete from minio: %w", err)
	}

	return nil
}
