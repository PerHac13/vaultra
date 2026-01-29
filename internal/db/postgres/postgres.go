package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/PerHac13/vaultra/internal/db"
)

type PostgreSQL struct {
	logger *slog.Logger
	config Config
}

func New(logger *slog.Logger, config Config) *PostgreSQL {
	return &PostgreSQL{
		logger: logger,
		config: config,
	}
}

func (p *PostgreSQL) Connect(ctx context.Context) error {
	p.logger.Debug("Connecting to PostgreSQL", "host", p.config.Host, "port", p.config.Port)

	cmd := exec.CommandContext(ctx, "psql",
		"-h", p.config.Host,
		"-p", fmt.Sprintf("%d", p.config.Port),
		"-U", p.config.User,
		"-d", p.config.Database,
		"-c", "SELECT 1",
	)	

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", p.config.Password))

	if err:= cmd.Run(); err != nil {
		p.logger.Error("Failed to connect to PostgreSQL", "error", err)
		return err
	}

	p.logger.Info("Successfully connected to PostgreSQL")
	return nil
}

func (p *PostgreSQL) Disconnect(ctx context.Context) error {
	p.logger.Info("Disconnecting from PostgreSQL")
	return nil
}

func (p *PostgreSQL) Ping(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "psql",
		"-h", p.config.Host,
		"-p", fmt.Sprintf("%d", p.config.Port),
		"-U", p.config.User,
		"-d", p.config.Database,
		"-c", "SELECT 1",
	)	

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", p.config.Password))

	return cmd.Run()
}

func (p *PostgreSQL) GetMetadata(ctx context.Context) (*db.Metadata, error) {
	return &db.Metadata{
		Name: p.config.Database,
		Version: "PostgreSQL",
	}, nil
}