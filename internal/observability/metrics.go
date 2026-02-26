package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics with label support
type Metrics struct {
	// Backup metrics - with db_type and storage_type labels
	BackupOperationsTotal   prometheus.CounterVec
	BackupSuccessTotal      prometheus.CounterVec
	BackupFailureTotal      prometheus.CounterVec
	BackupDurationSeconds   prometheus.GaugeVec
	BackupSizeBytes         prometheus.GaugeVec

	// Restore metrics - with db_type and storage_type labels
	RestoreOperationsTotal  prometheus.CounterVec
	RestoreSuccessTotal     prometheus.CounterVec
	RestoreFailureTotal     prometheus.CounterVec
	RestoreDurationSeconds  prometheus.GaugeVec

	// Storage metrics - with storage_type label
	StorageUploadBytesTotal      prometheus.CounterVec
	StorageUploadDurationSeconds prometheus.HistogramVec
	StorageDownloadBytesTotal    prometheus.CounterVec
	StorageDownloadDurationSeconds prometheus.HistogramVec

	// Database metrics - with db_type label
	DatabaseConnectionsTotal      prometheus.CounterVec
	DatabaseConnectionDurationSeconds prometheus.GaugeVec

	// Application metrics
	AppStartupTotal             prometheus.Counter
	AppUptimeSeconds            prometheus.Gauge
	AppConfigLoadDurationSeconds prometheus.Gauge
}

// NewMetrics creates and registers all metrics with proper labels
func NewMetrics() *Metrics {
	return &Metrics{
		// Backup metrics
		BackupOperationsTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "backup_operations_total",
				Help: "Total number of backup operations attempted",
			},
			[]string{"database_type", "storage_type"},
		),
		BackupSuccessTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "backup_success_total",
				Help: "Total number of successful backups",
			},
			[]string{"database_type", "storage_type"},
		),
		BackupFailureTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "backup_failure_total",
				Help: "Total number of failed backups",
			},
			[]string{"database_type", "storage_type"},
		),
		BackupDurationSeconds: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "backup_duration_seconds",
				Help: "Duration of last backup in seconds",
			},
			[]string{"database_type", "storage_type"},
		),
		BackupSizeBytes: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "backup_size_bytes",
				Help: "Size of last backup in bytes",
			},
			[]string{"database_type", "storage_type"},
		),

		// Restore metrics
		RestoreOperationsTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "restore_operations_total",
				Help: "Total number of restore operations attempted",
			},
			[]string{"database_type", "storage_type"},
		),
		RestoreSuccessTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "restore_success_total",
				Help: "Total number of successful restores",
			},
			[]string{"database_type", "storage_type"},
		),
		RestoreFailureTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "restore_failure_total",
				Help: "Total number of failed restores",
			},
			[]string{"database_type", "storage_type"},
		),
		RestoreDurationSeconds: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "restore_duration_seconds",
				Help: "Duration of last restore in seconds",
			},
			[]string{"database_type", "storage_type"},
		),

		// Storage metrics
		StorageUploadBytesTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "storage_upload_bytes_total",
				Help: "Total bytes uploaded to storage",
			},
			[]string{"storage_type"},
		),
		StorageUploadDurationSeconds: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "storage_upload_duration_seconds",
				Help:    "Duration of storage upload in seconds",
				Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
			},
			[]string{"storage_type"},
		),
		StorageDownloadBytesTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "storage_download_bytes_total",
				Help: "Total bytes downloaded from storage",
			},
			[]string{"storage_type"},
		),
		StorageDownloadDurationSeconds: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "storage_download_duration_seconds",
				Help:    "Duration of storage download in seconds",
				Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
			},
			[]string{"storage_type"},
		),

		// Database metrics
		DatabaseConnectionsTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_connections_total",
				Help: "Total number of database connections",
			},
			[]string{"database_type", "status"},
		),
		DatabaseConnectionDurationSeconds: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "database_connection_duration_seconds",
				Help: "Duration of last database connection in seconds",
			},
			[]string{"database_type"},
		),

		// Application metrics
		AppStartupTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "app_startup_total",
			Help: "Total number of application startups",
		}),
		AppUptimeSeconds: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "app_uptime_seconds",
			Help: "Application uptime in seconds",
		}),
		AppConfigLoadDurationSeconds: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "app_config_load_duration_seconds",
			Help: "Duration of configuration load in seconds",
		}),
	}
}