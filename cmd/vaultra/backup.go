package main

import (
	"context"
	"fmt"

	"github.com/PerHac13/vaultra/internal/app"
	"github.com/PerHac13/vaultra/internal/backup"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a database",
	Long:  `Create a backup of the specified database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if configFile == "" {
			return fmt.Errorf("--config flag is required")
		}
		ctx := context.Background()

		application, err := app.New(ctx, configFile)
		if err != nil {
			return fmt.Errorf("failed to initialize application: %w", err)
		}
		defer application.Close(ctx)

		req := backup.BackupRequest{
			Name:     "manual_backup",
			Strategy: backup.StrategyFull,
		}

		result, err := application.BackupEngine().Backup(ctx,req)
		if err != nil {
			return fmt.Errorf("backup failed: %w", err)
		}
		
		fmt.Printf("Backup completed successfully!\n")
		fmt.Printf("ID: %s\n", result.ID)
		fmt.Printf("Size: %d bytes\n", result.Size)
		fmt.Printf("Duration: %.2f seconds\n", result.Duration)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}