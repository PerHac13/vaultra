package compress

import (
	"context"
	"io"
)

type Compressor interface {
    Compress(ctx context.Context, r io.Reader, w io.Writer) error
    Decompress(ctx context.Context, r io.Reader, w io.Writer) error
	Name() string
}