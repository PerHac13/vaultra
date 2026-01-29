package postgres

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)


func (p *PostgreSQL) Restore(ctx context.Context, r io.Reader) error {
	p.logger.Info("Starting restore of PostgreSQL", "database", p.config.Database)

	cmd := exec.CommandContext(ctx, "psql",
		"-h", p.config.Host,
		"-p", fmt.Sprintf("%d", p.config.Port),
		"-U", p.config.User,
		"-d", p.config.Database,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", p.config.Password))
	cmd.Stdin = r
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		p.logger.Error("Failed to restore PostgreSQL", "error", err)
		return err
	}

	p.logger.Info("Successfully completed restore of PostgreSQL", "database", p.config.Database)
	return nil
}