package app

import (
	"context"
	"fmt"

	"github.com/PerHac13/vaultra/internal/backup"
	"github.com/PerHac13/vaultra/internal/config"
	"github.com/PerHac13/vaultra/internal/db"
	"github.com/PerHac13/vaultra/internal/db/mock"
	"github.com/PerHac13/vaultra/internal/db/postgres"
	"github.com/PerHac13/vaultra/internal/logging"
	"github.com/PerHac13/vaultra/internal/repository"
	"github.com/PerHac13/vaultra/internal/repository/inmemory"
	"github.com/PerHac13/vaultra/internal/restore"
	"github.com/PerHac13/vaultra/internal/storage"
	"github.com/PerHac13/vaultra/internal/storage/local"
)



type App struct {
	logger *logging.Logger
	config *config.ConfigType
	database db.Database
	storage  storage.Storage
	backupEngine *backup.Engine
	restoreEngine *restore.Engine
	repository repository.BackupRepository
}

func New(ctx context.Context, cfgFile string) (*App, error) {
	parser := config.NewParser()
	cfg, err := parser.ParseFile(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	logger := logging.NewDefaultLogger()
	validator := config.NewValidator(logger.Logger);

	if err := validator.Validator(cfg);  err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	logLevel := logging.ParseLevel(cfg.App.LogLevel)
	logger = logging.NewLogger(logLevel, nil)
	logger.Info("Configuration loaded", "name", cfg.App.Name)

	// Initialize database adapter based on config
	var database db.Database
	switch cfg.Database.Type {
	case "postgres":
		pgConfig := postgres.Config{
			Host:        getMapString(cfg.Database.Config, "host", "localhost"),
			Port:        getMapInt(cfg.Database.Config, "port", 5432),
			User:        getMapString(cfg.Database.Config, "user", "postgres"),
			Password:    getMapString(cfg.Database.Config, "password", ""),
			Database:    getMapString(cfg.Database.Config, "database", ""),
			SSLMode: 	 getMapString(cfg.Database.Config, "ssl_mode", "disable"),	
		}
		database = postgres.New(logger.Logger, pgConfig)
	case "mock":
		database = mock.NewMockDatabase(mock.ConfigType{
			Data: []byte("mock data"),
		})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	// Initialize storage adapter based on config
	var stor storage.Storage
	switch cfg.Storage.Type {
	case "local":
		basePath := getMapString(cfg.Storage.Config, "base_path", "./backups")
		s, err := local.NewLocalStorage(basePath)
		if err != nil {
			return nil, fmt.Errorf("initialize local storage: %w", err)
		}
		stor = s
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
	}

	// Initialize repo
	repo := inmemory.New()

	// Initialize backup and restore engines
	backupEngine := backup.New(logger.Logger, database, stor, repo)
	restoreEngine := restore.New(logger.Logger, database, stor, repo)

	return &App{
		logger: logger,
		config: cfg,
		database: database,
		storage:  stor,
		backupEngine: backupEngine,
		restoreEngine: restoreEngine,
		repository: repo,
	}, nil
}

func getMapString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

func getMapInt(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	return defaultValue
}

func (a *App) Logger() *logging.Logger {
	return a.logger
}

func (a *App) Config() *config.ConfigType {
	return a.config
}

func (a *App) BackupEngine() *backup.Engine {
	return a.backupEngine
}

func (a *App) RestoreEngine() *restore.Engine {
	return a.restoreEngine
}

func (a *App) BackupRepository() repository.BackupRepository {
	return a.repository;
}

func (a *App) Close(ctx context.Context) error {
	a.logger.Info("Shutting down application")
	return nil
}