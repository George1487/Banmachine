package minio

import (
	"context"
	"fmt"
	filestorage "ingestWorker/internal/storage/fileStorage"
	"io"

	"github.com/minio/minio-go/v7"
)

var _ filestorage.FileStorage = (*MiniOStorage)(nil)

type MiniOStorage struct {
	client *minio.Client
	bucket string
}

func NewMiniOStorage(client *minio.Client, bucket string) *MiniOStorage {
	return &MiniOStorage{
		client: client,
		bucket: bucket,
	}
}

func (s *MiniOStorage) GetFile(ctx context.Context, key string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object : %w", err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("read object : %w", err)
	}

	return data, nil
}

func (s *MiniOStorage) DeleteFile(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("delete object : %w", err)
	}
	return nil
}
