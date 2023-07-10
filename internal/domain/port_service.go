package domain

import (
	"fmt"
	"sync"
)

type PortService interface {
	CreateOrUpdatePorts() error
}

type portService struct {
	fileReader     FileRepository
	portRepository PortRepository
}

func NewPortService(fileRepository FileRepository, portRepository PortRepository) (PortService, error) {
	if fileRepository == nil {
		return nil, ErrNoFileRepository
	}

	if portRepository == nil {
		return nil, ErrNoPortRepository
	}

	return &portService{
		fileReader:     fileRepository,
		portRepository: portRepository,
	}, nil
}

func (s *portService) CreateOrUpdatePorts() error {
	portChan := make(chan Port)

	var wg sync.WaitGroup

	wg.Add(1)

	// Start a goroutine to read ports from the channel and save them to the database
	go func() {
		defer wg.Done()

		for port := range portChan {
			s.portRepository.UpsertPort(port)
		}
	}()

	err := s.fileReader.ReadAndReturnPorts(portChan)
	if err != nil {
		return fmt.Errorf("s.fileReader.ReadAndReturnPorts: %w", err)
	}

	wg.Wait()

	return nil
}
