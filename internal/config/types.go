package config

type ConfigType struct {
	App           AppConfig         `yaml:"app`
	Database      DatabaseConfig    `yaml:"database"`
	Storage       StorageConfig     `yaml:"storage"`
	Compression   CompressionConfig `yaml:"compression"`
	Backup        BackupConfig      `yaml:"backup"`
}

type AppConfig struct {
	Name      string      `yaml:"name"`
	LogLevel  string      `yaml:"log_level"`
}

type DatabaseConfig struct {
	Type    string        			`yaml:"type"`
	Config  map[string]interface{} 	`yaml:"config"`
}

type StorageConfig struct {
	Type    string        			`yaml:"type"`
	Config  map[string]interface{} 	`yaml:"config"`
}

type CompressionConfig struct {
	Algorithm   string  	 `yaml:"algorithm"`
	Level       int          `yaml:"level"`
}

type BackupConfig struct {
	RetentionDays  int    			`yaml:"retention_days"`
	FullBackup     FullBackupConfig `yaml:"full_backup"`
}

type FullBackupConfig struct {
	Enabled   bool   `yaml:"enabled"`
}