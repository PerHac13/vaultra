package mongodb

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func (m *MongoDB) Restore(ctx context.Context, r io.Reader) error {
	m.logger.Info("Starting restore of MongoDB database", "database", m.config.Database)

	cmd := exec.CommandContext(ctx, "mongorestore",
		"--host", m.config.Host,
		"--port", fmt.Sprintf("%d", m.config.Port),
		"--username", m.config.Username,
		"--password", m.config.Password,
		"--authenticationDatabase", m.config.AuthDatabase,
		"--db", m.config.Database,
		"--archive",
		"--drop",
	)

	cmd.Stdin = r
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		m.logger.Error("Failed to restore MongoDB database", "error", err)
		return err
	}

	m.logger.Info("Successfully completed restore of MongoDB database")

	return nil
}