package restore

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/PerHac13/vaultra/internal/db"
	"github.com/PerHac13/vaultra/internal/observability"
	"github.com/PerHac13/vaultra/internal/repository"
	"github.com/PerHac13/vaultra/internal/storage"
)

type Engine struct {
	logger       *slog.Logger
	db           db.Database
	storage      storage.Storage
	repo         repository.BackupRepository
	metrics      *observability.Metrics
	databaseType string
	storageType  string
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

type RestoreRequest struct {
	BackupPath string
	DryRun     bool
}

type RestoreResult struct {
	Duration  float64
	BytesRead int64
}

func (e *Engine) Restore(ctx context.Context, req RestoreRequest) (*RestoreResult, error) {
	startTime := time.Now()
	e.logger.Info("Starting restore process", "backup_path", req.BackupPath, "dry_run", req.DryRun)

	// Record operation start
	if e.metrics != nil {
		e.metrics.RestoreOperationsTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
	}

	// Connect to database
	dbConnectStart := time.Now()
	if err := e.db.Connect(ctx); err != nil {
		e.logger.Error("Failed to connect to database", "error", err)

		// Record failure metrics
		if e.metrics != nil {
			duration := time.Since(startTime).Seconds()
			e.metrics.RestoreFailureTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
			e.metrics.RestoreDurationSeconds.WithLabelValues(e.databaseType, e.storageType).Set(duration)
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

	// Download backup with metrics
	e.logger.Info("Downloading backup from storage", "path", req.BackupPath)

	downloadStartTime := time.Now()
	backupReader, err := e.storage.Download(ctx, req.BackupPath)
	if err != nil {
		e.logger.Error("Failed to download backup", "error", err)

		// Record failure metrics
		if e.metrics != nil {
			duration := time.Since(startTime).Seconds()
			e.metrics.RestoreFailureTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
			e.metrics.RestoreDurationSeconds.WithLabelValues(e.databaseType, e.storageType).Set(duration)
		}

		return nil, err
	}
	defer backupReader.Close()

	var bytesRead int64
	countingReader := &countingReader{
		reader: backupReader,
		count:  &bytesRead,
	}

	// Record storage download metrics
	downloadDuration := time.Since(downloadStartTime).Seconds()
	if e.metrics != nil {
		e.metrics.StorageDownloadBytesTotal.WithLabelValues(e.storageType).Add(float64(bytesRead))
		e.metrics.StorageDownloadDurationSeconds.WithLabelValues(e.storageType).Observe(downloadDuration)
	}

	if !req.DryRun {
		e.logger.Info("Restoring database from backup")
		restoreStartTime := time.Now()

		if err := e.db.Restore(ctx, countingReader); err != nil {
			e.logger.Error("Failed to restore database", "error", err)

			// Record failure metrics
			if e.metrics != nil {
				duration := time.Since(startTime).Seconds()
				e.metrics.RestoreFailureTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
				e.metrics.RestoreDurationSeconds.WithLabelValues(e.databaseType, e.storageType).Set(duration)
			}

			return nil, err
		}

		restoreDuration := time.Since(restoreStartTime).Seconds()
		if e.metrics != nil && restoreDuration > 0 {
			// Record restore duration separately if needed for deeper insights
			e.logger.Debug("Restore operation duration", "duration_seconds", restoreDuration)
		}
	} else {
		e.logger.Info("Dry run enabled, skipping actual restore")
	}

	duration := float64(time.Since(startTime).Milliseconds())

	// Record success metrics
	if e.metrics != nil {
		e.metrics.RestoreSuccessTotal.WithLabelValues(e.databaseType, e.storageType).Inc()
		e.metrics.RestoreDurationSeconds.WithLabelValues(e.databaseType, e.storageType).Set(duration / 1000) // Convert to seconds
	}

	e.logger.Info("Restore process completed",
		"duration_ms", duration,
		"bytes_read", bytesRead,
		"dry_run", req.DryRun,
	)

	return &RestoreResult{
		Duration:  duration,
		BytesRead: bytesRead,
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