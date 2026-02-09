package mongodb

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func (m *MongoDB) FullBackup(ctx context.Context, w io.Writer) error {
	m.logger.Info("Starting full backup of MongoDB database", "database", m.config.Database)

	cmd :=exec.CommandContext(ctx, "mongodump",
		"--host", m.config.Host,
		"--port", fmt.Sprintf("%d", m.config.Port),
		"--username", m.config.Username,
		"--password", m.config.Password,
		"--authenticationDatabase", m.config.AuthDatabase,
		"--db", m.config.Database,
		"--archive",
	)

	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		m.logger.Error("Failed to perform full backup of MongoDB database", "error", err)
		return err
	}
	
	m.logger.Info("Successfully completed full backup of MongoDB database")
	return nil
}

func (m *MongoDB) IncrementalBackup(ctx context.Context, w io.Writer, since time.Time) error {
	m.logger.Info("Incremental backups are not supported for MongoDB")
	return nil
}

