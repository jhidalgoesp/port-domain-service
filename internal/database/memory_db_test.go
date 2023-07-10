package database

import (
	"errors"
	"github.com/jhidalgoesp/ports/internal/domain"
	"strconv"
	"sync"
	"testing"
)

func TestDatabase_GetPortByID(t *testing.T) {
	testCases := []struct {
		name        string
		ports       map[string]*domain.Port
		id          string
		expected    *domain.Port
		expectedErr error
	}{
		{
			name: "Valid ID",
			ports: map[string]*domain.Port{
				"port1": {ID: "port1", Name: "Port 1"},
			},
			id:          "port1",
			expected:    &domain.Port{ID: "port1", Name: "Port 1"},
			expectedErr: nil,
		},
		{
			name:        "Invalid ID",
			ports:       map[string]*domain.Port{},
			id:          "port2",
			expected:    nil,
			expectedErr: ErrPortNotFound,
		},
	}

	db := NewDatabase()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Reset()

			for _, port := range tc.ports {
				db.UpsertPort(*port)
			}

			port, err := db.GetPortByID(tc.id)

			if tc.expectedErr != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if port.ID != tc.expected.ID {
					t.Errorf("Expected port: %v, got: %v", tc.expected, port)
				}
			}
		})
	}
}

func TestDatabase_UpsertPort(t *testing.T) {
	db := NewDatabase()

	port := domain.Port{ID: "port1", Name: "Port 1"}

	db.UpsertPort(port)

	got, err := db.GetPortByID(port.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if port.ID != got.ID {
		t.Errorf("Expected port: %v, got: %v", port, got)
	}
}

func TestDatabase_Reset(t *testing.T) {
	db := NewDatabase()

	port := domain.Port{ID: "port1", Name: "Port 1"}

	db.UpsertPort(port)

	db.Reset()

	got, err := db.GetPortByID(port.ID)
	if !errors.Is(err, ErrPortNotFound) {
		t.Errorf("Expected error: %v, got: %v", ErrPortNotFound, err)
	}

	if got != nil {
		t.Errorf("Expected nil port, got: %v", got)
	}
}

func TestDatabase_ConcurrentAccess(t *testing.T) {
	db := NewDatabase()

	numRoutines := 100
	wg := sync.WaitGroup{}
	wg.Add(numRoutines)

	for i := 0; i < numRoutines; i++ {
		go func(i int) {
			defer wg.Done()

			port := domain.Port{ID: "port", Name: "Port " + strconv.Itoa(i)}
			db.UpsertPort(port)

			got, err := db.GetPortByID("port")
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if got.ID != port.ID {
				t.Errorf("Expected port: %v, got: %v", port, got)
			}
		}(i)
	}

	wg.Wait()
}
