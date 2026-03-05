package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/PerHac13/vaultra/internal/app"
	"github.com/PerHac13/vaultra/internal/tui"
)

var (
	configFile string
	verbose    bool
	useTUI     bool
)

var rootCmd = &cobra.Command{
	Use:   "vaultra",
	Short: "Database backup and restore utility",
	Long:  "vaultra - Back up and restore databases reliably\n\nUse --tui for interactive mode or CLI subcommands for scripting.",
	Version: Version,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Allow help command without config
		if cmd.Name() == "help" {
			return nil
		}

		// TUI mode does not require config upfront
		if useTUI {
			return nil
		}

		// Require config for subcommands (backup, restore, list)
		if configFile == "" && cmd.Name() != "" && cmd.Use != "vaultra" {
			return fmt.Errorf("--config flag is required")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		if useTUI {
			runTUI(cmd.Context())
			return
		}

		// Otherwise, show help
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Persistent flags (available to all commands)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&useTUI, "tui", false, "Use interactive Terminal User Interface")
}

// runTUI launches the interactive Terminal User Interface
func runTUI(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	var vaultraApp *app.App
	var logger *slog.Logger

	// If config is provided, initialize app upfront
	if configFile != "" {
		var err error
		vaultraApp, err = app.New(ctx, configFile)
		if err != nil {
			log.Fatalf("Failed to initialize app: %v", err)
		}
		defer vaultraApp.Close(ctx)

		appLogger := vaultraApp.Logger()
		logger = appLogger.Logger
		logger.Info("Starting Vaultra in TUI mode", "config", configFile)
	} else {
		// No config — TUI will handle config setup
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
		logger.Info("Starting Vaultra in TUI mode (no config)")
	}

	// Create and launch TUI (vaultraApp may be nil)
	tuiApp := tui.New(vaultraApp, logger, configFile)
	if err := tuiApp.Run(); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}