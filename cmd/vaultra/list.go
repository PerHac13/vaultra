package main

import (
	"context"
	"fmt"

	"github.com/PerHac13/vaultra/internal/app"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backups",
	Long:  "Display all available backups.",
	RunE:  func(cmd *cobra.Command, args []string) error {
		if configFile == "" {
			return fmt.Errorf("--config flag is required")
		}

		ctx := context.Background()

		application, err := app.New(ctx, configFile)
		if err != nil {
			return fmt.Errorf("failed to initialize application: %w", err)
		}
		defer application.Close(ctx)

		backups, err := application.BackupRepository().List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list backups: %w", err)
		}

		if len(backups) == 0 {
			fmt.Println("No backups found.")
			return nil
		}

		fmt.Println("Available Backups:")
		fmt.Println("--------------------------------")
		for _, backup := range backups {
			fmt.Printf("ID: %s\n", backup.ID)
			fmt.Printf("Name: %s\n", backup.Name)
			fmt.Printf("Size: %d bytes\n", backup.Size)
			fmt.Printf("Created At: %s\n", backup.CreatedAt)
			fmt.Printf("Path: %s\n", backup.Path);
			fmt.Println("--------------------------------")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}