package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
    Upload(ctx context.Context, path string, data io.Reader) error
    Download(ctx context.Context, path string) (io.ReadCloser, error)
    List(ctx context.Context, prefix string) ([]BackupInfo, error)
    Delete(ctx context.Context, path string) error
}

type BackupInfo struct {
    Path        string
    Size        int64
   LastModified time.Time
}