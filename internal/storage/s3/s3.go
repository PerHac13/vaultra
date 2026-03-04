package s3

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/PerHac13/vaultra/internal/observability"
	"github.com/PerHac13/vaultra/internal/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	logger     *slog.Logger
	config     Config
	client     *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
	metrics    *observability.Metrics
}

func New(logger *slog.Logger, cfg Config) (*S3Storage, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("S3 bucket name is required")
	}

	// Set default values for optional fields
	if cfg.Region == "" {
		cfg.Region = DefaultRegion
	}
	if cfg.Prefix == "" {
		cfg.Prefix = DefaultPrefix
	}
	if cfg.PartSize <= 0 {
		cfg.PartSize = DefaultPartSize
	}

	logger.Debug("Initializing S3 storage",
		"bucket", cfg.Bucket,
		"region", cfg.Region,
		"prefix", cfg.Prefix,
	)

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region))

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = cfg.PartSize
	})

	downloader := manager.NewDownloader(client)

	return &S3Storage{
		logger:     logger,
		config:     cfg,
		client:     client,
		uploader:   uploader,
		downloader: downloader,
	}, nil
}

// SetMetrics sets the metrics instance for recording
func (s *S3Storage) SetMetrics(metrics *observability.Metrics) {
	s.metrics = metrics
}

func (s *S3Storage) Upload(ctx context.Context, path string, data io.Reader) error {
	startTime := time.Now()

	s.logger.Info("Uploading backup to S3", "bucket", s.config.Bucket, "path", path)

	fullPath := s.config.Prefix + path

	// Wrap reader to count bytes
	var bytesUploaded int64
	countingReader := &countingReader{
		reader: data,
		count:  &bytesUploaded,
	}

	result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(fullPath),
		Body:   countingReader,
	})
	if err != nil {
		s.logger.Error("Failed to upload to S3", "error", err)
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	duration := time.Since(startTime).Seconds()

	s.logger.Info("Upload successful",
		"location", result.Location,
		"bytes", bytesUploaded,
		"duration_seconds", duration,
	)

	// Record metrics
	if s.metrics != nil {
		s.metrics.StorageUploadBytesTotal.WithLabelValues("s3").Add(float64(bytesUploaded))
		s.metrics.StorageUploadDurationSeconds.WithLabelValues("s3").Observe(duration)
	}

	return nil
}

func (s *S3Storage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	startTime := time.Now()

	s.logger.Info("Downloading backup from S3", "bucket", s.config.Bucket, "path", path)

	fullPath := s.config.Prefix + path

	buf := manager.NewWriteAtBuffer([]byte{})
	n, err := s.downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(fullPath),
	})

	if err != nil {
		s.logger.Error("Failed to download from S3", "error", err)
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}

	duration := time.Since(startTime).Seconds()

	s.logger.Info("Download successful",
		"bytes_downloaded", n,
		"duration_seconds", duration,
	)

	// Record metrics
	if s.metrics != nil {
		s.metrics.StorageDownloadBytesTotal.WithLabelValues("s3").Add(float64(n))
		s.metrics.StorageDownloadDurationSeconds.WithLabelValues("s3").Observe(duration)
	}

	return io.NopCloser(strings.NewReader(string(buf.Bytes()))), nil
}

func (s *S3Storage) List(ctx context.Context, prefix string) ([]storage.BackupInfo, error) {
	s.logger.Info("Listing backups in S3", "bucket", s.config.Bucket, "prefix", prefix)

	fullPrefix := s.config.Prefix + prefix

	var backups []storage.BackupInfo
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.config.Bucket),
		Prefix: aws.String(fullPrefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			s.logger.Error("Failed to list objects in S3", "error", err)
			return nil, fmt.Errorf("failed to list objects in S3: %w", err)
		}

		for _, obj := range page.Contents {
			backups = append(backups, storage.BackupInfo{
				Path:         strings.TrimPrefix(aws.ToString(obj.Key), s.config.Prefix),
				Size:         aws.ToInt64(obj.Size),
				LastModified: aws.ToTime(obj.LastModified),
			})
		}
	}

	s.logger.Info("List successful", "backup_count", len(backups))
	return backups, nil
}

func (s *S3Storage) Delete(ctx context.Context, path string) error {
	s.logger.Info("Deleting backup from S3", "bucket", s.config.Bucket, "path", path)

	fullPath := s.config.Prefix + path

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(fullPath),
	})
	if err != nil {
		s.logger.Error("Failed to delete object from S3", "error", err)
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	s.logger.Info("Delete successful")
	return nil
}
type countingReader struct {
	reader io.Reader
	count  *int64
}

func (cr *countingReader) Read(p []byte) (int, error) {
	n, err := cr.reader.Read(p)
	*cr.count += int64(n)
	return n, err
}