package db_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/PerHac13/vaultra/internal/db/mock"
)

func TestMockConnect(t *testing.T) {
	db := mock.NewMockDatabase(mock.ConfigType{
		Data: []byte("test"),
	})

	err := db.Connect(context.Background())
	defer db.Disconnect(context.Background())
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}

	if !db.Connected {
		t.Fatalf("db should be connected")
	}
}

func TestMockFullBackup(t *testing.T) {
	testData := []byte("backup data")

	db := mock.NewMockDatabase(mock.ConfigType{
		Data: testData,
	})

	db.Connect(context.Background())
	defer db.Disconnect(context.Background())

	var buf bytes.Buffer
	err := db.FullBackup(context.Background(), &buf)
	if err != nil {
		t.Fatalf("full backup failed: %v", err)
	}

	if !bytes.Equal(buf.Bytes(), testData) {
		t.Fatalf("backup data does not match original data")
	}

}

func TestMockRestore(t *testing.T) {
	db := mock.NewMockDatabase(mock.ConfigType{Data: []byte("original")})
	db.Connect(context.Background())
	defer db.Disconnect(context.Background())

	newData := []byte("restored data")
	err := db.Restore(context.Background(), bytes.NewReader(newData))
	if err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	if string(db.Data) != string(newData) {
		t.Fatalf("data was not restored correctly")
	}
}

func TestMockFailAt(t *testing.T){
	db := mock.NewMockDatabase(mock.ConfigType{
		Data:  []byte("test"),
		FailAt: "Backeup",
	})
	db.Connect(context.Background())
	defer db.Disconnect(context.Background())
	
	var buf bytes.Buffer
	err := db.FullBackup(context.Background(), &buf)
	if err != nil {
		t.Fatalf("expected backup to fail, but it succeeded")
	}
}