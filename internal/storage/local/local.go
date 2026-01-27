package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/PerHac13/vaultra/internal/storage"
)

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("create base path: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (ls *LocalStorage) Upload(ctx context.Context, path string, data io.Reader) error {
	fullPath := filepath.Join(ls.basePath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("create directories: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, data)

	return err
}

func (ls *LocalStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(ls.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	return file, nil
}

func (ls *LocalStorage) List(ctx context.Context, prefix string) ([]storage.BackupInfo, error) {
	var infos []storage.BackupInfo
	searchPath := filepath.Join(ls.basePath, prefix)

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir(){
			return nil
		}

		relPath, _ := filepath.Rel(ls.basePath, path)
		infos = append(infos, storage.BackupInfo{
			Path:       relPath,
			Size:       info.Size(),
			LastModified: info.ModTime(),
		})
		return nil
	})

	return infos, err
}

func (ls *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(ls.basePath, path)

	return os.Remove(fullPath)
}