package db

import (
	"context"
	"io"
	"time"
)

type Database interface {
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    Ping(ctx context.Context) error
    FullBackup(ctx context.Context, w io.Writer) error
    IncrementalBackup(ctx context.Context, w io.Writer, since time.Time) error
    Restore(ctx context.Context, r io.Reader) error
    GetMetadata(ctx context.Context) (*Metadata, error)
}

type Metadata struct {
    Name       string
    Size       int64
    Version    string
}

type ErrorKind string

const (
    KindConnection     ErrorKind = "ConnectionError"
    KindDatabase        ErrorKind = "DatabaseError"
)

type Error struct {
    Err     error
    Kind    ErrorKind
}

func (e *Error) Error() string {
    return e.Err.Error()
}

func (e *Error) Unwrap() error {
    return e.Err
}