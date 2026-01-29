package main

import (
	"context"
	"fmt"

	"github.com/PerHac13/vaultra/internal/app"
	"github.com/PerHac13/vaultra/internal/restore"
	"github.com/spf13/cobra"
)

var (
	backupPath string
	dryRun     bool
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a backup",
	Long:  "Restore a backup from storage to the database",
	RunE:  func(cmd *cobra.Command, args []string) error {
		if configFile == "" {
			return fmt.Errorf("--config flag is required")
		}
		if backupPath == "" {
			return fmt.Errorf("--backup-path flag is required")
		}

		ctx := context.Background()

		application, err := app.New(ctx, configFile)
		if err != nil {
			return fmt.Errorf("failed to initialize application: %w", err)
		}
		defer application.Close(ctx)

		req := restore.RestoreRequest{
			BackupPath: backupPath,
			DryRun:     dryRun,
		}

		result, err := application.RestoreEngine().Restore(ctx, req)
		if err != nil {
			return fmt.Errorf("restore failed: %w", err)
		}
		
		fmt.Printf("Restore completed successfully. Duration: %.2f seconds\n", result.Duration)
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringVarP(&backupPath, "backup-path", "p", "", "Path to the backup to restore")
	restoreCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Perform a dry run without making changes")

}