package app

import (
	"context"
	"fmt"

	"github.com/PerHac13/vaultra/internal/config"
	"github.com/PerHac13/vaultra/internal/logging"
)



type App struct {
	logger *logging.Logger
	config *config.ConfigType
}

func New(ctx context.Context, cfgFile string) (*App, error) {
	parser := config.NewParser()
	cfg, err := parser.ParseFile(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	logger := logging.NewDefaultLogger()
	validator := config.NewValidator(logger.Logger);

	if err := validator.Validator(cfg);  err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	logLevel := logging.ParseLevel(cfg.App.LogLevel)
	logger = logging.NewLogger(logLevel, nil)
	logger.Info("Configuration loaded", "name", cfg.App.Name)

	return &App{
		logger: logger,
		config: cfg,
	}, nil
}

func (a *App) Logger() *logging.Logger {
	return a.logger
}

func (a *App) Config() *config.ConfigType {
	return a.config
}

func (a *App) Close(ctx context.Context) error {
	a.logger.Info("Shutting down application")
	return nil
}