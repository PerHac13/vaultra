package restore

import (
	"context"
	"log/slog"
	"time"

	"github.com/PerHac13/vaultra/internal/db"
	"github.com/PerHac13/vaultra/internal/repository"
	"github.com/PerHac13/vaultra/internal/storage"
)

type Engine struct {
	logger  *slog.Logger
	db      db.Database
	storage storage.Storage
	repo    repository.BackupRepository
}

func New(
	logger *slog.Logger,
	database db.Database,
	stor storage.Storage,
	repo repository.BackupRepository,
) *Engine {
	return &Engine{
		logger:  logger,
		db:      database,
		storage: stor,
		repo:    repo,
	}
}

type RestoreRequest struct {
	BackupPath    string
	DryRun        bool
}

type RestoreResult struct {
	Duration float64
}

func (e *Engine) Restore(ctx context.Context, req RestoreRequest) (*RestoreResult, error) {
	startTime := time.Now()
	e.logger.Info("Starting restore process", "backup_path", req.BackupPath, "dry_run", req.DryRun)

	// Connect to database
	if err := e.db.Connect(ctx); err != nil {
		e.logger.Error("Failed to connect to database", "error", err)
		return nil, err
	}
	defer e.db.Disconnect(ctx)

	// Download backup 
	e.logger.Info("Downloading backup from storage", "path", req.BackupPath)

	backupReader, err := e.storage.Download(ctx, req.BackupPath)
	if err != nil {
		e.logger.Error("Failed to download backup", "error", err)
		return nil, err
	}
	defer backupReader.Close()

	if !req.DryRun {
		e.logger.Info("Restoring database from backup")
		if err := e.db.Restore(ctx, backupReader); err != nil {
			e.logger.Error("Failed to restore database", "error", err)
			return nil, err
		}
	} else {
		e.logger.Info("Dry run enabled, skipping actual restore")
	} 

	duration := float64(time.Since(startTime).Milliseconds())
	e.logger.Info("Restore process completed", "duration_ms", duration)

	return &RestoreResult{
		Duration: duration,
	}, nil
}