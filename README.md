# Vaultra

**Vaultra** is a universal, database-agnostic backup and restore utility written in Go.
It provides a unified CLI and TUI to back up, restore, and manage databases across multiple database engines with support for compression, cloud storage, scheduling, and full observability using OpenTelemetry.

This project follows the **Standard Go Project Layout** as a guideline to ensure maintainability, scalability, and clarity as the codebase grows.

---

## Overview

Modern systems often rely on multiple database technologies across environments. Managing backups consistently becomes complex and error-prone. vaultra addresses this by providing:

* A single interface for multiple DBMS
* Streaming, efficient backups for large databases
* Pluggable storage backends
* First-class observability
* Cross-platform support

vaultra is designed for developers, SREs, and platform engineers who want reliable backups without vendor lock-in.

---

## Key Features

### Database Support

* MySQL
* PostgreSQL
* MongoDB
* SQLite
  (Designed to be extensible to other databases)

### Backup Capabilities

* Full backups
* Incremental backups (where supported)
* Differential backups
* Streaming backups (no intermediate dump files)

### Restore Capabilities

* Full restore
* Selective restore (tables, schemas, collections where supported)
* Dry-run restore validation

### Storage Backends

* Local filesystem
* AWS S3
* Google Cloud Storage
* Azure Blob Storage

### Compression

* gzip
* zstd (optional)

### Observability

* Structured logs via OpenTelemetry
* Metrics for duration, size, and failures
* Traces across backup and restore pipelines

### Interfaces

* CLI for automation and scripting
* TUI for interactive usage and monitoring

---

## Project Goals

* Provide a database-agnostic backup and restore tool
* Minimize performance impact on database servers
* Handle large databases efficiently
* Offer production-grade logging and diagnostics
* Remain portable across Linux, macOS, and Windows

### Non-Goals

* Real-time replication or CDC
* Database high-availability management
* Replacing native database HA/DR solutions

---

## Repository Structure

This repository follows the **Standard Go Project Layout** as a guideline.

```
vaultra/
├── cmd/                # Application entry points (CLI/TUI)
│   └── vaultra/
│       └── main.go
├── internal/           # Private application code
│   ├── app/            # Application orchestration
│   ├── backup/         # Backup engine
│   ├── restore/        # Restore engine
│   ├── db/             # Database adapters
│   ├── storage/        # Storage backends
│   ├── compress/       # Compression logic
│   ├── scheduler/      # Backup scheduling
│   ├── observability/  # OpenTelemetry setup
│   └── tui/            # TUI implementation
├── pkg/                # Public, reusable libraries (if any)
├── configs/            # Configuration templates
├── scripts/            # Build, lint, and utility scripts
├── build/              # Dockerfiles, packaging, CI configs
├── deployments/        # Docker Compose, Kubernetes, etc.
├── docs/               # Design docs and architecture notes
├── examples/           # Usage examples
├── tools/              # Supporting developer tools
├── test/               # Integration tests and test data
├── go.mod
├── go.sum
└── README.md
```

The `internal` directory is used to enforce encapsulation and prevent unintended reuse of private packages.

---

## CLI Usage (Planned)

```bash
vaultra backup --db postgres --config ./configs/postgres.yaml
vaultra restore --from s3://backups/prod/db.tar.zst
vaultra list backups
vaultra doctor
```

---

## TUI

The TUI provides:

* Interactive database configuration
* Backup and restore progress visualization
* History and logs viewer
* Safe confirmation flows for destructive actions

The TUI is optional and does not replace the CLI.

---

## Observability

vaultra uses OpenTelemetry for observability:

* Logs: structured, machine-readable
* Metrics: backup duration, size, error rates
* Traces: end-to-end visibility into backup and restore workflows

Compatible with Prometheus, Grafana, Loki, and Jaeger.

---

## Security Considerations

* Credentials via environment variables or config files
* No plaintext secrets in logs
* TLS for database connections where supported
* Optional encryption at rest for backups
* Least-privilege cloud IAM policies

---

## Tooling

* Language: Go
* CLI: Cobra (planned)
* TUI: Bubble Tea
* Observability: OpenTelemetry
* Environment management: mise
* Containerization: Docker
* CI: GitHub Actions (planned)

---

## Development Setup

```bash
mise install
mise use
go build ./cmd/vaultra
```

Docker:

```bash
docker build -t vaultra .
```

---

## Roadmap

### v0.1

* CLI
* MySQL and PostgreSQL support
* Local storage
* Full backups

### v0.2

* MongoDB support
* Cloud storage (S3)
* Compression

### v0.3

* TUI
* Incremental backups
* Notifications

### v1.0

* Full cloud support
* Metrics dashboards
* Restore validation and safety checks

---

