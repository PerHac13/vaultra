package local

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/PerHac13/vaultra/internal/observability"
	"github.com/PerHac13/vaultra/internal/storage"
)



func New(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("create base path: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func NewWithMetrics(basePath string, logger *slog.Logger, metrics *observability.Metrics) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("create base path: %w", err)
	}
	return &LocalStorage{
		basePath: basePath,
		logger:   logger,
		metrics:  metrics,
	}, nil
}

func (ls *LocalStorage) Upload(ctx context.Context, path string, data io.Reader) error {
	startTime := time.Now()

	fullPath := filepath.Join(ls.basePath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("create directories: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	// Copy data and count bytes
	bytesWritten, err := io.Copy(file, data)
	if err != nil {
		return fmt.Errorf("write data: %w", err)
	}

	duration := time.Since(startTime).Seconds()

	if ls.logger != nil {
		ls.logger.Debug("Uploaded to local storage",
			"path", path,
			"bytes", bytesWritten,
			"duration_seconds", duration,
		)
	}

	// Record metrics
	if ls.metrics != nil {
		ls.metrics.StorageUploadBytesTotal.WithLabelValues("local").Add(float64(bytesWritten))
		ls.metrics.StorageUploadDurationSeconds.WithLabelValues("local").Observe(duration)
	}

	return nil
}

func (ls *LocalStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	startTime := time.Now()

	fullPath := filepath.Join(ls.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	// Get file size for metrics
	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("stat file: %w", err)
	}

	size := fileInfo.Size()
	duration := time.Since(startTime).Seconds()

	if ls.logger != nil {
		ls.logger.Debug("Downloaded from local storage",
			"path", path,
			"bytes", size,
			"duration_seconds", duration,
		)
	}

	// Record metrics
	if ls.metrics != nil {
		ls.metrics.StorageDownloadBytesTotal.WithLabelValues("local").Add(float64(size))
		ls.metrics.StorageDownloadDurationSeconds.WithLabelValues("local").Observe(duration)
	}

	return file, nil
}

func (ls *LocalStorage) List(ctx context.Context, prefix string) ([]storage.BackupInfo, error) {
	var infos []storage.BackupInfo
	searchPath := filepath.Join(ls.basePath, prefix)

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(ls.basePath, path)
		infos = append(infos, storage.BackupInfo{
			Path:         relPath,
			Size:         info.Size(),
			LastModified: info.ModTime(),
		})
		return nil
	})

	return infos, err
}

func (ls *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(ls.basePath, path)

	if ls.logger != nil {
		ls.logger.Debug("Deleting from local storage", "path", path)
	}

	return os.Remove(fullPath)
}

// SetMetrics sets the metrics instance (for compatibility with S3 pattern)
func (ls *LocalStorage) SetMetrics(metrics *observability.Metrics) {
	ls.metrics = metrics
}

// SetLogger sets the logger instance
func (ls *LocalStorage) SetLogger(logger *slog.Logger) {
	ls.logger = logger
}