package observability

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsServer manages the Prometheus metrics HTTP server
type MetricsServer struct {
	port   int
	logger *slog.Logger
	server *http.Server
	mu     sync.Mutex
}

// NewMetricsServer creates a new metrics server
func NewMetricsServer(logger *slog.Logger, port int) *MetricsServer {
	return &MetricsServer{
		port:   port,
		logger: logger,
	}
}

// Start starts the metrics HTTP server
func (ms *MetricsServer) Start(ctx context.Context) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%d", ms.port)
	ms.logger.Info("Starting metrics server", "port", ms.port, "endpoint", fmt.Sprintf("http://localhost:%d/metrics", ms.port))

	ms.server = &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		if err := ms.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ms.logger.Error("Metrics server error", "error", err)
		}
	}()

	// Wait for context or server stop
	go func() {
		<-ctx.Done()
		ms.Stop()
	}()

	return nil
}

// Stop gracefully stops the metrics server
func (ms *MetricsServer) Stop() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.server == nil {
		return nil
	}

	ms.logger.Info("Stopping metrics server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ms.server.Shutdown(ctx)
}

// RecordBackupMetrics records metrics for a backup operation
func RecordBackupMetrics(metrics *Metrics, dbType, storageType, compression string, duration time.Duration, size int64, success bool) {
	metrics.BackupOperationsTotal.WithLabelValues(dbType, storageType).Inc()

	if success {
		metrics.BackupSuccessTotal.WithLabelValues(dbType, storageType).Inc()
		metrics.BackupDurationSeconds.WithLabelValues(dbType, storageType).Set(duration.Seconds())
		metrics.BackupSizeBytes.WithLabelValues(dbType, storageType).Set(float64(size))
	} else {
		metrics.BackupFailureTotal.WithLabelValues(dbType, storageType).Inc()
	}
}

// RecordRestoreMetrics records metrics for a restore operation
func RecordRestoreMetrics(metrics *Metrics, dbType, storageType string, duration time.Duration, success bool) {
	metrics.RestoreOperationsTotal.WithLabelValues(dbType, storageType).Inc()

	if success {
		metrics.RestoreSuccessTotal.WithLabelValues(dbType, storageType).Inc()
		metrics.RestoreDurationSeconds.WithLabelValues(dbType, storageType).Set(duration.Seconds())
	} else {
		metrics.RestoreFailureTotal.WithLabelValues(dbType, storageType).Inc()
	}
}

// RecordStorageUploadMetrics records metrics for storage upload
func RecordStorageUploadMetrics(metrics *Metrics, storageType string, bytes int64, duration time.Duration) {
	metrics.StorageUploadBytesTotal.WithLabelValues(storageType).Add(float64(bytes))
	metrics.StorageUploadDurationSeconds.WithLabelValues(storageType).Observe(duration.Seconds())
}

// RecordStorageDownloadMetrics records metrics for storage download
func RecordStorageDownloadMetrics(metrics *Metrics, storageType string, bytes int64, duration time.Duration) {
	metrics.StorageDownloadBytesTotal.WithLabelValues(storageType).Add(float64(bytes))
	metrics.StorageDownloadDurationSeconds.WithLabelValues(storageType).Observe(duration.Seconds())
}

// RecordDatabaseConnectionMetrics records metrics for database connections
func RecordDatabaseConnectionMetrics(metrics *Metrics, dbType string, duration time.Duration, success bool) {
	if success {
		metrics.DatabaseConnectionsTotal.WithLabelValues(dbType, "success").Inc()
		metrics.DatabaseConnectionDurationSeconds.WithLabelValues(dbType).Set(duration.Seconds())
	} else {
		metrics.DatabaseConnectionsTotal.WithLabelValues(dbType, "failure").Inc()
	}
}

// RecordConfigLoadMetrics records metrics for config loading
func RecordConfigLoadMetrics(metrics *Metrics, duration time.Duration) {
	metrics.AppStartupTotal.Inc()
	metrics.AppConfigLoadDurationSeconds.Set(duration.Seconds())
}