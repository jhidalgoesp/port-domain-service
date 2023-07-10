package file_reader

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jhidalgoesp/ports/internal/domain"
	"log"
	"os"
)

type FileReader struct {
	FilePath string
}

func NewFileReader(filePath string) (*FileReader, error) {
	if filePath == "" {
		return nil, ErrNoFilePath
	}

	return &FileReader{
		FilePath: filePath,
	}, nil
}

func (f *FileReader) ReadAndReturnPorts(portChan chan<- domain.Port) error {
	file, err := os.Open(f.FilePath)

	log.Println(f.FilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrFileDoesNotExist
		}

		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	// Read the opening curly brace
	_, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("error reading JSON: %w", err)
	}

	// Iterate over each key-value pair
	for decoder.More() {
		// Read the port ID
		key, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("error reading JSON: %w", err)
		}

		var port domain.Port

		err = decoder.Decode(&port)
		if err != nil {
			return fmt.Errorf("error decoding JSON: %w", err)
		}

		port.ID = key.(string)

		portChan <- port
	}

	// Read the closing curly brace
	_, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("error reading JSON: %w", err)
	}

	close(portChan)

	return nil
}
