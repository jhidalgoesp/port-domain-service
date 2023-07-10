package database

import (
	"fmt"
	"github.com/jhidalgoesp/ports/internal/domain"
	"sync"
)

type MemDB struct {
	ports map[string]*domain.Port
	mu    sync.RWMutex
}

func NewDatabase() *MemDB {
	return &MemDB{
		ports: make(map[string]*domain.Port),
	}
}

func (db *MemDB) GetPortByID(id string) (*domain.Port, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	port, ok := db.ports[id]
	if !ok {
		return nil, ErrPortNotFound
	}

	return port, nil
}

func (db *MemDB) UpsertPort(port domain.Port) {
	db.mu.Lock()
	defer db.mu.Unlock()

	fmt.Printf("Saved port: %v\n", port)

	db.ports[port.ID] = &port
}

func (db *MemDB) Reset() {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.ports = make(map[string]*domain.Port)
}
