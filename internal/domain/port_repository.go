package domain

type PortRepository interface {
	GetPortByID(id string) (*Port, error)
	UpsertPort(port Port)
}
