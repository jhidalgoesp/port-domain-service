package domain

type FileRepository interface {
	ReadAndReturnPorts(portChan chan<- Port) error
}
