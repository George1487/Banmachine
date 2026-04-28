package filestorage

import (
	"context"
)

type FileStorage interface {
	GetFile(ctx context.Context, key string) ([]byte, error)
	DeleteFile(ctx context.Context, key string) error
}
