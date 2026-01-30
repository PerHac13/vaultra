# Vaultra: Design

## High-Level Architecture (Preliminary)

```
User Input (CLI/TUI)
        ↓
   [App Orchestrator]
        ↓
    ┌───┴───────┬──────────┬────────────┐
    ↓           ↓          ↓            ↓
[Backup]   [Restore]  [List/Info]  [Schedule]
    ↓           ↓          ↓            ↓
    └───┬───────┴──────────┴────────────┘
        ↓
  [Database Adapter]  ← Abstract interface
  (PostgreSQL, MySQL, etc.)
        ↓
  [Storage Backend]   ← Abstract interface
  (Local, S3, GCS, etc.)
        ↓
  [Compression]       ← Abstract interface
  (gzip, zstd, etc.)
        ↓
  [Observability]     ← Logging, metrics, traces
```

---

## Architecture Layers

```
┌──────────────────────────────────────────────────┐
│          CLI Interface (Cobra)                    │
│     backup | restore | list | schedule | doctor  │
└──────────────────────┬───────────────────────────┘
                       │
┌──────────────────────┴───────────────────────────┐
│     Application Orchestrator (App Container)     │
│     - Dependency Injection                       │
│     - Configuration Management                   │
│     - Error Handling                             │
└──────────────────────┬───────────────────────────┘
                       │
┌──────────┬───────────┴───────────┬───────────────┐
│          │                       │               │
v          v                       v               v
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ Backup       │ │ Restore      │ │ List/Info    │
│ Engine       │ │ Engine       │ │ Commands     │
└──────────────┘ └──────────────┘ └──────────────┘
       │                  │              │
┌──────┴──────────────────┼──────────────┴───────┐
│                         │                      │
v                         v                      v
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ Database     │ │ Storage      │ │ Compression  │
│ Adapter      │ │ Backend      │ │ Strategy     │
│ (Interface)  │ │ (Interface)  │ │ (Interface)  │
└──────────────┘ └──────────────┘ └──────────────┘
       │                │                │
┌──────┴──┬──────┐  ┌───┴───┬──────┐  ┌──┴──┬──────┐
│          │      │  │       │      │  │     │      │
v          v      v  v       v      v  v     v      v
Postgres MySQL Mongo Local S3 GCS Azure Gzip Zstd None
```
