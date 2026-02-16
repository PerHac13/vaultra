package s3

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

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

func (s *S3Storage) Upload(ctx context.Context, path string, data io.Reader) error {
	s.logger.Info("Uploading backup to S3","bucker", s.config.Bucket, "path", path)

	fullPath := s.config.Prefix + path

	result , err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &s.config.Bucket,
		Key:    &fullPath,
		Body:   data,
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}
	s.logger.Info("Upload successful", "location", result.Location)
	return nil
}

func (s *S3Storage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	s.logger.Info("Downloading backup from S3", "bucket", s.config.Bucket, "path", path)

	fullPath := s.config.Prefix + path

	buf := manager.NewWriteAtBuffer([]byte{})
	n, err := s.downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: &s.config.Bucket,
		Key:    &fullPath,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}
	s.logger.Info("Download successful", "bytes_downloaded", n)
	return io.NopCloser(strings.NewReader(string(buf.Bytes()))), nil
}

func (s *S3Storage) List(ctx context.Context, prefix string) ([]storage.BackupInfo, error) {
	s.logger.Info("Listing backups in S3", "bucket", s.config.Bucket, "prefix", prefix)

	fullPrefix := s.config.Prefix + prefix

	var backups []storage.BackupInfo
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: &s.config.Bucket,
		Prefix: &fullPrefix,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects in S3: %w", err)
		}

		for _, obj := range page.Contents {
			backups = append(backups, storage.BackupInfo{
				Path: strings.TrimPrefix(*obj.Key, s.config.Prefix),
				Size: *obj.Size,
				LastModified: *obj.LastModified,
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
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}
	s.logger.Info("Delete successful")
	return nil
}