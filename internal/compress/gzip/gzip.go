package gzip

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
)

type GzipCompressor struct{
	level int
}

func New(level int)  (*GzipCompressor, error) {
	if level < 1 || level > 9 {
		level = 6
	}
	return &GzipCompressor{level: level}, nil
}

func (gc *GzipCompressor) Name() string {
	return "gzip"
}

func (gc *GzipCompressor) Compress(ctx context.Context, r io.Reader, w io.Writer) error {
	gw, err := gzip.NewWriterLevel(w, gc.level)

	if err != nil {
		return fmt.Errorf("create gzip writer: %w", err)
	}
	defer gw.Close()

	_, err = io.Copy(gw, r)
	return err
}

func (gc *GzipCompressor) Decompress(ctx context.Context, r io.Reader, w io.Writer) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("create gzip reader: %w", err)
	}
	defer gr.Close()

	_, err = io.Copy(w, gr)
	return err
}

