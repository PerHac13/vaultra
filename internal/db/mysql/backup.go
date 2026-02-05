package mysql

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func (m *MySQL) FullBackup(ctx context.Context, w io.Writer) error {
	m.logger.Info("Starting full backup of MySQL database", "database", m.config.Database)

	cmd := exec.CommandContext(ctx, "mysqldump",
		"-h", m.config.Host,
		"-P", fmt.Sprintf("%d", m.config.Port),
		"-u", m.config.User,
		fmt.Sprintf("-p%s", m.config.Password),
		"--single-transaction",
		"--lock-tables=false",
		"--set-charset",
		"--default-character-set", m.config.Charset,
		m.config.Database,
	)

	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		m.logger.Error("Failed to perform full backup of MySQL database", "error", err)
		return err
	}

	m.logger.Info("Successfully completed full backup of MySQL database")
	return nil
}

func (m *MySQL) IncrementalBackup(ctx context.Context, w io.Writer, since time.Time) error {
	m.logger.Info("Incremental backups are not supported for MySQL")
	return nil
}