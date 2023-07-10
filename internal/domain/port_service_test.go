package domain

import (
	"errors"
	"fmt"
	"testing"
)

// MockFileRepository is a mock implementation of the FileRepository interface
type MockFileRepository struct {
	ReadAndReturnPortsFunc func(portChan chan<- Port) error
}

func (m *MockFileRepository) ReadAndReturnPorts(portChan chan<- Port) error {
	if m.ReadAndReturnPortsFunc != nil {
		return m.ReadAndReturnPortsFunc(portChan)
	}

	return nil
}

// MockPortRepository is a mock implementation of the PortRepository interface
type MockPortRepository struct {
	GetPortByIDFunc func(id string) (*Port, error)
	UpsertPortFunc  func(port Port)
}

func (m *MockPortRepository) GetPortByID(id string) (*Port, error) {
	if m.GetPortByIDFunc != nil {
		return m.GetPortByIDFunc(id)
	}

	return nil, nil
}

func (m *MockPortRepository) UpsertPort(port Port) {
	if m.UpsertPortFunc != nil {
		m.UpsertPortFunc(port)
	}
}

func TestNewPortService(t *testing.T) {
	testCases := []struct {
		name           string
		fileRepository FileRepository
		portRepository PortRepository
		expectedError  error
	}{
		{
			name:           "Success",
			fileRepository: &MockFileRepository{},
			portRepository: &MockPortRepository{},
			expectedError:  nil,
		},
		{
			name:           "Nil file repository",
			fileRepository: nil,
			portRepository: &MockPortRepository{},
			expectedError:  ErrNoFileRepository,
		},
		{
			name:           "Nil port repository",
			fileRepository: &MockFileRepository{},
			portRepository: nil,
			expectedError:  ErrNoPortRepository,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewPortService(tc.fileRepository, tc.portRepository)

			if err != nil && !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestPortService_CreateOrUpdatePorts(t *testing.T) {
	testCases := []struct {
		name             string
		fileReadError    error
		expectedError    error
		expectedPortIDs  []string
		expectedNumCalls int
	}{
		{
			name:             "Successful read",
			fileReadError:    nil,
			expectedError:    nil,
			expectedPortIDs:  []string{"AEAJM", "AEAUH", "AEDXB"},
			expectedNumCalls: 3,
		},
		{
			name:             "File read error",
			fileReadError:    errors.New("file read error"),
			expectedError:    fmt.Errorf("s.fileReader.ReadAndReturnPorts: file read error"),
			expectedPortIDs:  nil,
			expectedNumCalls: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock objects
			fileRepository := &MockFileRepository{}
			portRepository := &MockPortRepository{}

			// Create test instance of PortService
			service, err := NewPortService(fileRepository, portRepository)
			if err != nil {
				t.Errorf("Error creating PortService: %v", err)
			}

			// Prepare test data
			ports := make([]Port, len(tc.expectedPortIDs))
			for i, portID := range tc.expectedPortIDs {
				ports[i] = Port{ID: portID, Name: "Port " + portID}
			}

			// Set expectations for the file repository
			fileRepository.ReadAndReturnPortsFunc = func(portChan chan<- Port) error {
				for _, port := range ports {
					portChan <- port
				}
				close(portChan)
				return tc.fileReadError
			}

			// Set expectations for the port repository
			var numCalls int
			portRepository.UpsertPortFunc = func(port Port) {
				numCalls++
			}

			// Call the method being tested
			err = service.CreateOrUpdatePorts()

			// Assert that the error matches the expected error
			if tc.fileReadError != nil {
				if err == nil || errors.Is(err, tc.expectedError) {
					t.Errorf("Expected error: %v, got: %v", tc.fileReadError, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Assert that the expected number of calls were made
			if numCalls != tc.expectedNumCalls {
				t.Errorf("Expected %d calls to UpsertPort, got %d", tc.expectedNumCalls, numCalls)
			}
		})
	}
}
