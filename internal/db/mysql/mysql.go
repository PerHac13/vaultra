package mysql

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/PerHac13/vaultra/internal/db"
)

type MySQL struct {
	logger   *slog.Logger
	config   Config
}

func New(logger *slog.Logger, config Config) *MySQL {
	if config.Charset == "" {
		config.Charset = DefaultCharset
	}
	if config.Port == 0 {
		config.Port = DefaultPort
	}
	return &MySQL{
		logger: logger,
		config: config,
	}
}

func (m *MySQL) Connect(ctx context.Context) error {
	m.logger.Info("Connecting to MySQL", "host", m.config.Host, "port", m.config.Port, "database", m.config.Database)

	cmd := exec.CommandContext(ctx, "mysql",
		"-h", m.config.Host,
		"-P", fmt.Sprintf("%d", m.config.Port),
		"-u", m.config.User,
		fmt.Sprintf("-p%s", m.config.Password),
		"-e", "SELECT 1",
		m.config.Database,
	)

	if err := cmd.Run(); err != nil {
		m.logger.Error("Failed to connect to MySQL", "error", err)
		return err
	}

	m.logger.Info("Successfully connected to MySQL")
	return nil
}

func (m *MySQL) Disconnect(ctx context.Context) error {
	m.logger.Info("Disconnecting from MySQL")
	return nil
}

func (m *MySQL) Ping(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "mysql",
		"-h", m.config.Host,
		"-P", fmt.Sprintf("%d", m.config.Port),
		"-u", m.config.User,
		fmt.Sprintf("-p%s", m.config.Password),
		"-e", "SELECT 1",
		m.config.Database,
	)

	if err := cmd.Run(); err != nil {
		m.logger.Error("Failed to ping MySQL", "error", err)
		return err
	}

	m.logger.Info("Successfully pinged MySQL")
	return nil
}

func (m *MySQL) GetMetadata(ctx context.Context) (*db.Metadata, error) {
	return &db.Metadata{
		Name: m.config.Database,
		Version: "MySQL",
	}, nil
}