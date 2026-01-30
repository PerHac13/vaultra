package config

import (
	"fmt"
	"log/slog"
)


type Validator struct {
	logger *slog.Logger
}

func NewValidator(logger *slog.Logger) *Validator {
	return &Validator{
		logger: logger,
	}
}

func (v *Validator) Validator(cfg *ConfigType) error {
	var errs []string

	if cfg.App.Name == "" {
		errs = append(errs, "App.Name cannot be empty")
	}

	if cfg.App.LogLevel == "" {
		errs = append(errs, "App.LogLevel cannot be empty")
	}

	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLevels[cfg.App.LogLevel] {
		errs = append(errs, "App.LogLevel must be one of: debug, info, warn, error")
	}

	if err := v.validateDatabase(cfg.Database); err != nil {
		errs = append(errs, fmt.Sprintf("database: %v",err))
	}

	if err := v.validateStorage(cfg.Storage); err != nil {
		errs = append(errs, fmt.Sprintf("storage: %v",err))
	}

	if err:= v.validateCompression(cfg.Compression); err != nil {
		errs = append(errs, fmt.Sprintf("compression: %v",err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("config validation failed:\n%v", errs)
	}

	return nil
}

func (v *Validator) validateDatabase(cfg DatabaseConfig) error {
	if cfg.Type == "" {
		return fmt.Errorf("Type cannot be empty")
	}

	validTypes := map[string]bool{
		"mysql":    true,
		"postgres": true,
		"sqlite":   true,
		"mongodb":  true,
	}

	if !validTypes[cfg.Type] {
		return fmt.Errorf("Type must be one of: mysql, postgres, sqlite, mongodb")
	}

	if cfg.Config == nil {
		return fmt.Errorf("Config cannot be empty")
	}

	return nil
}

func (v *Validator) validateStorage(cfg StorageConfig) error {
	if cfg.Type == "" {
		return fmt.Errorf("Type cannot be empty")
	}

	validTypes := map[string]bool{
		"local":    true,
		"s3":       true,
		"gcs":      true,
		"azure":    true,
	}

	if !validTypes[cfg.Type] {
		return fmt.Errorf("Type must be one of: local, s3, gcs, azure")
	}

	if cfg.Config == nil {
		return fmt.Errorf("Config cannot be empty")
	}

	return nil
}

func (v *Validator) validateCompression(cfg CompressionConfig) error {
	if cfg.Algorithm == "" {
		return fmt.Errorf("Algorithm cannot be empty")
	}

	validAlgorithms := map[string]bool{
		"gzip":  true,
		"lz4":   true,
		"zstd":  true,
		"none":  true,
	}

	if !validAlgorithms[cfg.Algorithm] {
		return fmt.Errorf("Algorithm must be one of: gzip, lz4, zstd, snappy")
	}

	if cfg.Algorithm != "none" && (cfg.Level < 1 || cfg.Level > 9) {
		return fmt.Errorf("Level must be between 1 and 9 for selected Algorithm")
	}

	return nil
}	