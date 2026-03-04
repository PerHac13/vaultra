package local

import (
	"log/slog"

	"github.com/PerHac13/vaultra/internal/observability"
)

type LocalStorage struct {
	basePath string
	logger   *slog.Logger
	metrics  *observability.Metrics
}