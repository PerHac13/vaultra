package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var (
	configFile string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "vaultra",
	Short: "Database backup and restore utility",
	Long:  "vaultra - Back up and restore databases reliably",
	Version: Version,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if configFile == "" && cmd.Use != "vaultra" {
			return fmt.Errorf("--config flag is required")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
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
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}