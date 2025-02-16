package storage

type Storage interface {
	PutObject() error
	GetObject() error
}
