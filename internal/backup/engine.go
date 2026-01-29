package backup

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/PerHac13/vaultra/internal/db"
	"github.com/PerHac13/vaultra/internal/repository"
	"github.com/PerHac13/vaultra/internal/storage"
)

type Engine struct {
	logger        *slog.Logger
	db            db.Database
	storage       storage.Storage
	compressor    interface{}
	repo          repository.BackupRepository
}

func New(
	logger *slog.Logger,
	database db.Database,
	stor storage.Storage,
	repo repository.BackupRepository,
) *Engine {
	return &Engine{
		logger:     logger,
		db:         database,
		storage:    stor,
		repo:       repo,
	}
}

type BackupRequest struct {
	Name	 string
	Strategy Strategy
}

type BackupResult struct {
	ID       string
	Size     int64
	Duration float64
}

func (e *Engine) CreateBackup(ctx context.Context, req BackupRequest) (*BackupResult, error) {
	startTime := time.Now()

	e.logger.Info("Starting backup", "name", req.Name, "strategy", req.Strategy)


	metadata := &Metadata{
		ID:		fmt.Sprintf("backup_%d", startTime.Unix()),
		Name:	   req.Name,
		StartTime: startTime,
		Status: "in_progress",
	}


	// Connect to the database
	if err := e.db.Connect(ctx); err != nil {
		metadata.Status = "failed"
		metadata.Error  = err.Error()
		e.logger.Error("Database connection failed", "error", err)
		return nil, err
	}
	defer e.db.Disconnect(ctx)

	// Ping the database
	if err := e.db.Ping(ctx); err != nil {
		metadata.Status = "failed"
		metadata.Error  = err.Error()
		e.logger.Error("Database ping failed", "error", err)
		return nil, err
	}

	reader, writer := io.Pipe()
	defer reader.Close()

	errChan := make(chan error, 1)

	go func() {
		defer writer.Close()

		switch req.Strategy {
		case StrategyFull:
			errChan <- e.db.FullBackup(ctx, writer)
		default:
			errChan <- fmt.Errorf("unsupported backup strategy: %s", req.Strategy)
		}
	}()

	// Upload backup
	backupPath := fmt.Sprintf("backups/%s_%d.sql", req.Name, startTime.Unix())

	if err := e.storage.Upload(ctx, backupPath, reader); err != nil {
		metadata.Status = "failed"
		metadata.Error  = err.Error()
		e.logger.Error("Backup upload failed", "error", err)
		return nil, err
	}


	// Check for errors from the backup goroutine
	if err := <-errChan; err != nil {
		metadata.Status = "failed"
		metadata.Error  = err.Error()
		e.logger.Error("Backup creation failed", "error", err)
		return nil, err
	}

	endTime := time.Now()
	// metadata.EndTime = endTime
	metadata.Duration = float64(endTime.Sub(startTime).Milliseconds())
	metadata.Status = "success"

	backupRepo := &repository.Backup{
		ID:         metadata.ID,
		Name:       metadata.Name,
		Size:       0,  
		CreatedAt:  metadata.StartTime,
		Path:       backupPath,
		Status:     metadata.Status,
	}

	if err := e.repo.Save(ctx, backupRepo); err != nil {
		e.logger.Error("Failed to save backup metadata", "error", err)
	}
	e.logger.Info("Backup completed successfully", "id", metadata.ID, "duration_ms", metadata.Duration)

	return &BackupResult{
		ID:       metadata.ID,
		Duration: metadata.Duration,
	}, nil
}