package app

import (
	"context"
	"fmt"
	"time"

	"github.com/PerHac13/vaultra/internal/backup"
	"github.com/PerHac13/vaultra/internal/config"
	"github.com/PerHac13/vaultra/internal/db"
	"github.com/PerHac13/vaultra/internal/db/mock"
	"github.com/PerHac13/vaultra/internal/db/mongodb"
	"github.com/PerHac13/vaultra/internal/db/mysql"
	"github.com/PerHac13/vaultra/internal/db/postgres"
	"github.com/PerHac13/vaultra/internal/logging"
	"github.com/PerHac13/vaultra/internal/observability"
	"github.com/PerHac13/vaultra/internal/repository"
	"github.com/PerHac13/vaultra/internal/repository/inmemory"
	"github.com/PerHac13/vaultra/internal/restore"
	"github.com/PerHac13/vaultra/internal/storage"
	"github.com/PerHac13/vaultra/internal/storage/local"
	"github.com/PerHac13/vaultra/internal/storage/s3"
)



type App struct {
	logger *logging.Logger
	config *config.ConfigType
	database db.Database
	storage  storage.Storage
	backupEngine *backup.Engine
	restoreEngine *restore.Engine
	repository repository.BackupRepository
	metricsServer *observability.MetricsServer
	metrics       *observability.Metrics
	startTime      time.Time
	databaseType   string
	storageType	   string
}

func New(ctx context.Context, cfgFile string) (*App, error) {
	startTime := time.Now()

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

	// Initialize metrics and its server
	metrics := observability.NewMetrics()
	observability.RecordConfigLoadMetrics(metrics, time.Since(startTime))

	metricsPort := getMapInt(cfg.Observability.Config, "metrics_port", 9090)
	metricsServer := observability.NewMetricsServer(logger.Logger, metricsPort)

	if err := metricsServer.Start(ctx); err != nil {
		return nil, fmt.Errorf("start metrics server: %w", err)
	}

	// Initialize database adapter based on config
	var database db.Database
	databaseType := cfg.Database.Type
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
	case "mqsql":
		mysqlConfig := mysql.Config{
			Host:     getMapString(cfg.Database.Config, "host", "localhost"),
			Port:     getMapInt(cfg.Database.Config, "port", 3306),
			User:     getMapString(cfg.Database.Config, "user", "root"),
			Password: getMapString(cfg.Database.Config, "password", ""),
			Database: getMapString(cfg.Database.Config, "database", ""),
			Charset:  getMapString(cfg.Database.Config, "charset", "utf8mb4"),
    	}
		database = mysql.New(logger.Logger, mysqlConfig)
	case "mongodb":
		mongodbConfig := mongodb.Config{
			Host:         getMapString(cfg.Database.Config, "host", "localhost"),
			Port:         getMapInt(cfg.Database.Config, "port", 27017),
			Username:     getMapString(cfg.Database.Config, "username", ""),
			Password:     getMapString(cfg.Database.Config, "password", ""),
			Database:     getMapString(cfg.Database.Config, "database", ""),
			AuthDatabase: getMapString(cfg.Database.Config, "auth_database", "admin"),
			URI:          getMapString(cfg.Database.Config, "uri", ""),
		}
		database = mongodb.New(logger.Logger, mongodbConfig)
	case "mock":
		database = mock.NewMockDatabase(mock.ConfigType{
			Data: []byte("mock data"),
		})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	// Initialize storage adapter based on config
	var stor storage.Storage
	storageType := cfg.Storage.Type
	switch cfg.Storage.Type {
	case "local":
		basePath := getMapString(cfg.Storage.Config, "base_path", "./backups")
		s, err := local.New(basePath)
		if err != nil {
			return nil, fmt.Errorf("initialize local storage: %w", err)
		}
		stor = s
	case "s3":
		s3Config := s3.Config{
			Bucket:          getMapString(cfg.Storage.Config, "bucket", ""),
			Region:          getMapString(cfg.Storage.Config, "region", "us-east-1"),
			AccessKeyID:     getMapString(cfg.Storage.Config, "access_key_id", ""),
			SecretAccessKey: getMapString(cfg.Storage.Config, "secret_access_key", ""),
			Prefix:          getMapString(cfg.Storage.Config, "prefix", "vaultra/backups/"),
			Endpoint:        getMapString(cfg.Storage.Config, "endpoint", ""),
			DisableSSL:      getMapBool(cfg.Storage.Config, "disable_ssl", false),
		}
		s3Storage, err := s3.New(logger.Logger, s3Config)
		if err != nil {
			return nil, fmt.Errorf("create s3 storage: %w", err)
		}
		stor = s3Storage
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
	}

	// Initialize repo
	repo := inmemory.New()

	// Initialize backup and restore engines
	backupEngine := backup.New(logger.Logger, database, stor, repo, metrics, databaseType, storageType)
	restoreEngine := restore.New(logger.Logger, database, stor, repo, metrics, databaseType, storageType)


	app := &App{
		logger: logger,
		config: cfg,
		database: database,
		storage:  stor,
		backupEngine: backupEngine,
		restoreEngine: restoreEngine,
		repository: repo,
		metricsServer: metricsServer,
		metrics: metrics,
		startTime: startTime,
		databaseType: databaseType,
		storageType: storageType,
	}

	// Start recording uptime metric
	go app.recordUptime()

	return app, nil
}

func (a *App) recordUptime() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		a.metrics.AppUptimeSeconds.Set(time.Since(a.startTime).Seconds())
	} 
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

func (a *App) Metrics() *observability.Metrics {
	return a.metrics
}

func (a *App) Close(ctx context.Context) error {
	a.logger.Info("Shutting down application")
	if a.metricsServer != nil {
		if err := a.metricsServer.Stop(); err != nil {
			a.logger.Error("Failed to stop metrics server", "error", err)
		}
	}
	return nil
}