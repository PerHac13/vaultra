package mongodb

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/PerHac13/vaultra/internal/db"
)

type MongoDB struct {
	logger   *slog.Logger
	config   Config
}

func New(logger *slog.Logger, config Config) *MongoDB {
	if config.Port == 0 {
		config.Port = DefaultPort
	}
	if config.AuthDatabase == "" {
		config.AuthDatabase = DefaultAuthDatabase
	}
	return &MongoDB{
		logger: logger,
		config: config,
	}
}

func (m *MongoDB) Connect( ctx context.Context) error {
	m.logger.Debug("Connecting to MongoDB", "host", m.config.Host, "port", m.config.Port, "database", m.config.Database)

	cmd := exec.CommandContext(ctx, "mongosh",
		"--host", m.config.Host,
		"--port", fmt.Sprintf("%d", m.config.Port),
		"--username", m.config.Username,
		"--password", m.config.Password,
		"--authenticationDatabase", m.config.AuthDatabase,
		"--eval", "db.adminCommand('ping')",
		// m.config.Database,
	)

	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Fallback to mongo command if mongosh is not available
		m.logger.Warn("mongosh command failed, falling back to mongo", "error", err)
		cmd = exec.CommandContext(ctx, "mongo",
			"--host", m.config.Host,
			"--port", fmt.Sprintf("%d", m.config.Port),
			"--username", m.config.Username,
			"--password", m.config.Password,
			"--authenticationDatabase", m.config.AuthDatabase,
			"--eval", "db.adminCommand('ping')",
			// m.config.Database,
		)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to connect to MongoDB: %w", err)
		}
	}
	
	m.logger.Info("Successfully connected to MongoDB")
	return nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	m.logger.Info("Disconnected from MongoDB")
	return nil
}

func (m *MongoDB) Ping(ctx context.Context) error {
	m.logger.Debug("Pinging MongoDB", "host", m.config.Host, "port", m.config.Port, "database", m.config.Database)

	cmd := exec.CommandContext(ctx, "mongosh",
		"--host", m.config.Host,
		"--port", fmt.Sprintf("%d", m.config.Port),
		"--username", m.config.Username,
		"--password", m.config.Password,
		"--authenticationDatabase", m.config.AuthDatabase,
		"--eval", "db.adminCommand('ping')",
	)

	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		m.logger.Warn("mongosh command failed, falling back to mongo", "error", err)
		cmd = exec.CommandContext(ctx, "mongo",
			"--host", m.config.Host,
			"--port", fmt.Sprintf("%d", m.config.Port),
			"--username", m.config.Username,
			"--password", m.config.Password,
			"--authenticationDatabase", m.config.AuthDatabase,
			"--eval", "db.adminCommand('ping')",
		)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to ping MongoDB: %w", err)
		}
	}

	m.logger.Info("Successfully pinged MongoDB")
	return nil
}

func (m *MongoDB) GetMetadata(ctx context.Context) (*db.Metadata, error) {
	return &db.Metadata{
		Name: m.config.Database,
		Version: "MongoDB",
	}, nil
}