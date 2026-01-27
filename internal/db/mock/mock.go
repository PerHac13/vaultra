package mock

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/PerHac13/vaultra/internal/db"
)


type MockDatabase struct {
	Data       []byte
	Connected  bool
	FailAt     string
	Metadata   *db.Metadata
}


type ConfigType struct {
	Data       []byte
	FailAt	   string
	Metadata   *db.Metadata
}

func NewMockDatabase(cfg ConfigType) *MockDatabase {
	if cfg.Metadata == nil {
		cfg.Metadata = &db.Metadata{
			Name:    "MockDB",
			Size:    int64(len(cfg.Data)),
			Version: "1.0.0-mock",
		}
	}
	return &MockDatabase{
		Data:      cfg.Data,
		Connected: false,
		FailAt:    cfg.FailAt,
		Metadata:  cfg.Metadata,
	}
}

func (m *MockDatabase) Connect(ctx context.Context) error {
	if m.FailAt == "Connect" {
		return fmt.Errorf("connect failed")
	}

	m.Connected = true
	return nil
}

func (m *MockDatabase) Disconnect(ctx context.Context) error {
	if m.FailAt == "Disconnect" {
		return fmt.Errorf("disconnect failed")
	}

	m.Connected = false
	return nil
}

func (m *MockDatabase) Ping(ctx context.Context) error {
	if m.FailAt == "Ping" {
		return fmt.Errorf("ping failed")
	}
	if !m.Connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (m *MockDatabase) FullBackup(ctx context.Context, w io.Writer) error {
	if m.FailAt == "backup" {
		return fmt.Errorf("backup failed")
	}
	if !m.Connected {
		return fmt.Errorf("not connected")
	}
	_, err := w.Write(m.Data)
	return err
}

func (m *MockDatabase) IncrementalBackup(ctx context.Context, w io.Writer, since time.Time) error {
	return m.FullBackup(ctx, w)
}

func (m *MockDatabase) Restore(ctx context.Context, r io.Reader) error {
	if m.FailAt == "restore" {
		return fmt.Errorf("restore failed")
	}
	if !m.Connected {
		return fmt.Errorf("not connected")
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	m.Data = data
	return nil
}

func (m *MockDatabase) GetMetadata(ctx context.Context) (*db.Metadata, error){
	if !m.Connected {
		return nil, fmt.Errorf("not connected")
	}
	m.Metadata.Size = int64(len(m.Data))
	return m.Metadata, nil
}
