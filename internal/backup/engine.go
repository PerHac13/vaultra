package backup

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/PerHac13/vaultra/internal/db"
	"github.com/PerHac13/vaultra/internal/observability"
	"github.com/PerHac13/vaultra/internal/repository"
	"github.com/PerHac13/vaultra/internal/storage"
)

type Engine struct {
	logger        *slog.Logger
	db            db.Database
	storage       storage.Storage
	repo          repository.BackupRepository
	metrics       *observability.Metrics
	databaseType  string
	storageType   string
}

func New(
	logger *slog.Logger,
	database db.Database,
	stor storage.Storage,
	repo repository.BackupRepository,
	metrics *observability.Metrics,
	databaseType string,
	storageType string,
) *Engine {
	return &Engine{
		logger:       logger,
		db:           database,
		storage:      stor,
		repo:         repo,
		metrics:      metrics,
		databaseType: databaseType,
		storageType:  storageType,
	}
}

type BackupRequest struct {
	Name     string
	Strategy Strategy
}

type BackupResult struct {
	ID       string
	Size     int64
	Duration float64
}

func (e *Engine) Backup(ctx context.Context, req BackupRequest) (*BackupResult, error) {
	startTime := time.Now()

	e.logger.Info("Starting backup", "name", req.Name, "strategy", req.Strategy)

	// Record operation start
	if e.metrics != nil {
		e.metrics.BackupOperationsTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
	}

	metadata := &Metadata{
		ID:        fmt.Sprintf("backup_%d", startTime.Unix()),
		Name:      req.Name,
		StartTime: startTime,
		Status:    "in_progress",
	}

	// Connect to the database
	dbConnectStart := time.Now()
	if err := e.db.Connect(ctx); err != nil {
		metadata.Status = "failed"
		metadata.Error = err.Error()
		e.logger.Error("Database connection failed", "error", err)

		// Record failure metrics
		if e.metrics != nil {
			duration := time.Since(startTime).Seconds()
			e.metrics.BackupFailureTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
			e.metrics.BackupDurationSeconds.WithLabelValues(e.databaseType, e.storageType).Set(duration)
			e.metrics.DatabaseConnectionsTotal.WithLabelValues(e.databaseType, "failure").Inc()
		}

		return nil, err
	}
	defer e.db.Disconnect(ctx)

	// Record successful connection
	if e.metrics != nil {
		connDuration := time.Since(dbConnectStart).Seconds()
		e.metrics.DatabaseConnectionsTotal.WithLabelValues(e.databaseType, "success").Inc()
		e.metrics.DatabaseConnectionDurationSeconds.WithLabelValues(e.databaseType).Set(connDuration)
	}

	// Ping the database
	if err := e.db.Ping(ctx); err != nil {
		metadata.Status = "failed"
		metadata.Error = err.Error()
		e.logger.Error("Database ping failed", "error", err)

		// Record failure metrics
		if e.metrics != nil {
			duration := time.Since(startTime).Seconds()
			e.metrics.BackupFailureTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
			e.metrics.BackupDurationSeconds.WithLabelValues(e.databaseType, e.storageType).Set(duration)
		}

		return nil, err
	}

	reader, writer := io.Pipe()
	defer reader.Close()

	errChan := make(chan error, 1)
	var backupSize int64

	countingReader := &countingReader{
		reader: reader,
		count:  &backupSize,
	}

	go func() {
		defer writer.Close()

		switch req.Strategy {
		case StrategyFull:
			errChan <- e.db.FullBackup(ctx, writer)
		default:
			errChan <- fmt.Errorf("unsupported backup strategy: %s", req.Strategy)
		}
	}()

	// Upload backup with metrics
	backupPath := fmt.Sprintf("backups/%s_%d.sql", req.Name, startTime.Unix())
	uploadStartTime := time.Now()

	if err := e.storage.Upload(ctx, backupPath, countingReader); err != nil {
		metadata.Status = "failed"
		metadata.Error = err.Error()
		e.logger.Error("Backup upload failed", "error", err)

		// Record failure metrics
		if e.metrics != nil {
			duration := time.Since(startTime).Seconds()
			e.metrics.BackupFailureTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
			e.metrics.BackupDurationSeconds.WithLabelValues(e.databaseType, e.storageType).Set(duration)
		}

		return nil, err
	}

	// Record storage upload metrics
	uploadDuration := time.Since(uploadStartTime).Seconds()
	if e.metrics != nil {
		e.metrics.StorageUploadBytesTotal.WithLabelValues(e.storageType).Add(float64(backupSize))
		e.metrics.StorageUploadDurationSeconds.WithLabelValues(e.storageType).Observe(uploadDuration)
	}

	// Check for errors from the backup goroutine
	if err := <-errChan; err != nil {
		metadata.Status = "failed"
		metadata.Error = err.Error()
		e.logger.Error("Backup creation failed", "error", err)

		// Record failure metrics
		if e.metrics != nil {
			duration := time.Since(startTime).Seconds()
			e.metrics.BackupFailureTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
			e.metrics.BackupDurationSeconds.WithLabelValues(e.databaseType, e.storageType).Set(duration)
		}

		return nil, err
	}

	endTime := time.Now()
	metadata.Duration = float64(endTime.Sub(startTime).Milliseconds())
	metadata.Status = "success"

	backupRepo := &repository.Backup{
		ID:        metadata.ID,
		Name:      metadata.Name,
		Size:      backupSize,
		CreatedAt: metadata.StartTime,
		Path:      backupPath,
		Status:    metadata.Status,
	}

	if err := e.repo.Save(ctx, backupRepo); err != nil {
		e.logger.Error("Failed to save backup metadata", "error", err)
	}

	// Record success metrics
	if e.metrics != nil {
		e.metrics.BackupSuccessTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
		e.metrics.BackupDurationSeconds.WithLabelValues(e.databaseType, e.storageType).Set(metadata.Duration / 1000) // Convert to seconds
		e.metrics.BackupSizeBytes.WithLabelValues(e.databaseType, e.storageType).Set(float64(backupSize))
	}

	e.logger.Info("Backup completed successfully",
		"id", metadata.ID,
		"duration_ms", metadata.Duration,
		"size_bytes", backupSize,
		"storage_type", e.storageType,
	)

	return &BackupResult{
		ID:       metadata.ID,
		Size:     backupSize,
		Duration: metadata.Duration,
	}, nil
}

type countingReader struct {
	reader io.ReadCloser
	count  *int64
}

func (cr *countingReader) Read(p []byte) (int, error) {
	n, err := cr.reader.Read(p)
	*cr.count += int64(n)
	return n, err
}

func (cr *countingReader) Close() error {
	return cr.reader.Close()
}