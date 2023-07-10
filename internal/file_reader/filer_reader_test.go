package file_reader

import (
	"encoding/json"
	"errors"
	"github.com/jhidalgoesp/ports/internal/domain"
	"os"
	"sync"
	"testing"
)

func TestNewPortService(t *testing.T) {
	testCases := []struct {
		name          string
		filePath      string
		expectedError error
	}{
		{
			name:          "Valid repositories",
			filePath:      "test.json",
			expectedError: nil,
		},
		{
			name:          "Empty filepath",
			filePath:      "",
			expectedError: ErrNoFilePath,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewFileReader(tc.filePath)

			if err != nil && !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestFileReader_ReadAndReturnPorts(t *testing.T) {
	testCases := []struct {
		Name           string
		CreateTestFile bool
		ExpectedErr    error
		ExpectedPorts  []domain.Port
	}{
		{
			Name:           "Existing File",
			CreateTestFile: true,
			ExpectedErr:    nil,
			ExpectedPorts: []domain.Port{
				{
					ID:          "AEAJM",
					Name:        "Ajman",
					City:        "Ajman",
					Country:     "United Arab Emirates",
					Alias:       []string{},
					Regions:     []string{},
					Coordinates: []float64{55.5136433, 25.4052165},
					Province:    "Ajman",
					Timezone:    "Asia/Dubai",
					Unlocs:      []string{"AEAJM"},
					Code:        "52000",
				},
				{
					ID:          "AEAUH",
					Name:        "Abu Dhabi",
					City:        "Abu Dhabi",
					Country:     "United Arab Emirates",
					Alias:       []string{},
					Regions:     []string{},
					Coordinates: []float64{54.37, 24.47},
					Province:    "Abu Z¸aby [Abu Dhabi]",
					Timezone:    "Asia/Dubai",
					Unlocs:      []string{"AEAUH"},
					Code:        "52001",
				},
				{
					ID:          "AEDXB",
					Name:        "Dubai",
					City:        "Dubai",
					Country:     "United Arab Emirates",
					Alias:       []string{},
					Regions:     []string{},
					Coordinates: []float64{55.27, 25.25},
					Province:    "Dubayy [Dubai]",
					Timezone:    "Asia/Dubai",
					Unlocs:      []string{"AEDXB"},
					Code:        "52005",
				},
			},
		},
		{
			Name:           "Non-Existent File",
			CreateTestFile: false,
			ExpectedErr:    ErrFileDoesNotExist,
			ExpectedPorts:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fileName := "test.json"

			if tc.CreateTestFile {
				// Create a temporary file with the expected ports
				file, err := os.CreateTemp("", "test.json")
				if err != nil {
					t.Fatalf("Failed to create temporary file: %v", err)
				}
				defer os.Remove(file.Name())

				// Create a map to hold the expected ports
				portsMap := make(map[string]domain.Port)
				for _, port := range tc.ExpectedPorts {
					portsMap[port.ID] = port
				}

				// Write the ports map to the temporary file
				encoder := json.NewEncoder(file)
				err = encoder.Encode(portsMap)
				if err != nil {
					t.Fatalf("Failed to write expected ports to file: %v", err)
				}

				// Close the file
				file.Close()

				// Update the file path to the temporary file
				fileName = file.Name()
			}

			// Create a FileReader instance
			fileReader, err := NewFileReader(fileName)
			if err != nil {
				t.Fatalf("Failed to create FileReader: %v", err)
			}

			// Create a channel to receive the ports
			portChan := make(chan domain.Port)

			receivedPorts := make([]domain.Port, 0)

			var wg sync.WaitGroup

			wg.Add(1)

			// Start reading the ports in a separate goroutine
			go func() {
				defer wg.Done()

				for port := range portChan {
					receivedPorts = append(receivedPorts, port)
				}
			}()

			err = fileReader.ReadAndReturnPorts(portChan)
			if tc.ExpectedErr == nil {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}
			} else {
				if errors.Is(err, tc.ExpectedErr) {
					return
				}
			}

			wg.Wait()

			// Ensure the correct number of ports is received
			if len(receivedPorts) != len(tc.ExpectedPorts) {
				t.Fatalf("Received %d ports, expected %d", len(receivedPorts), len(tc.ExpectedPorts))
			}

			// Ensure the received ports match the expected ports
			for _, expectedPort := range tc.ExpectedPorts {
				found := false
				for _, receivedPort := range receivedPorts {
					if expectedPort.ID == receivedPort.ID {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("Expected port %v not found in received ports", expectedPort)
				}
			}

			if len(portChan) > 0 {
				t.Error("Port channel is not closed")
			}
		})
	}
}
