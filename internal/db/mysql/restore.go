package mysql

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func (m *MySQL) Restore(ctx context.Context, r io.Reader) error {
	m.logger.Info("Starting restore of MySQL database", "database", m.config.Database)

	cmd := exec.CommandContext(ctx, "mysql",
		"-h", m.config.Host,
		"-P", fmt.Sprintf("%d", m.config.Port),
		"-u", m.config.User,
		fmt.Sprintf("-p%s", m.config.Password),
		m.config.Database,
	)

	cmd.Stdin = r
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		m.logger.Error("Failed to restore MySQL database", "error", err)
		return err
	}

	m.logger.Info("Successfully completed restore of MySQL database")
	return nil
}