package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/PerHac13/vaultra/internal/repository"
)


type InMemoryRepository struct {
	mu             sync.RWMutex
	backups        map[string]*repository.Backup
}

func New() *InMemoryRepository {
	return &InMemoryRepository{
		backups: make(map[string]*repository.Backup),
	}
}

func (r *InMemoryRepository) Save(ctx context.Context, backup *repository.Backup) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if backup.ID == "" {
		return fmt.Errorf("backup id required")
	}

	r.backups[backup.ID] = backup
	
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*repository.Backup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	backup, exists := r.backups[id]
	if !exists {
		return nil, fmt.Errorf("backup not found: %s", id)
	}

	return backup, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]repository.Backup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var backups []repository.Backup
	for _, backup := range r.backups {
		backups = append(backups, *backup)
	}

	return backups, nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.backups, id)
	return nil
}