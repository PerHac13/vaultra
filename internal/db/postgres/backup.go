package postgres

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)


func (p *PostgreSQL) FullBackup(ctx context.Context, w io.Writer) error {
	p.logger.Info("Starting full backup of PostgreSQL", "database", p.config.Database)

	cmd := exec.CommandContext(ctx, "pg_dump",
		"-h", p.config.Host,
		"-p", fmt.Sprintf("%d", p.config.Port),
		"-U", p.config.User,
		p.config.Database,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", p.config.Password))
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		p.logger.Error("Failed to perform full backup of PostgreSQL", "error", err)
		return err
	}

	p.logger.Info("Successfully completed full backup of PostgreSQL", "database", p.config.Database)
	return nil
}

func (p *PostgreSQL) IncrementalBackup(ctx context.Context, w io.Writer, since time.Time) error {
	p.logger.Info("Incremental backups are not supported for PostgreSQL")
	return nil
}