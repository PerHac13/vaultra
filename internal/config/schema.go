package config

const SchemaVersion = "1.0.0"

func DefaultConfigSchema() *ConfigType {
	return &ConfigType{
		App: AppConfig{
			Name:     "vaultra",
			LogLevel: "info",
		},
		Compression: CompressionConfig{
			Algorithm: "gzip",
			Level: 6,
		},
		Backup: BackupConfig{
			RetentionDays: 30,
			FullBackup: FullBackupConfig{
				Enabled: true,
			},
		},
	}
}