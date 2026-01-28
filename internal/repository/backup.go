package repository

import (
	"context"
	"time"
)

type BackupRepository interface {
    Save(ctx context.Context, backup *Backup) error
    Get(ctx context.Context, id string) (*Backup, error)
    List(ctx context.Context) ([]Backup, error)
    Delete(ctx context.Context, id string) error
}

type Backup struct {
    ID        string
    Name      string
    Size      int64
    CreatedAt time.Time
    Path      string
    Status    string
}