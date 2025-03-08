package storage

import "io"

// Store is the interface for the object store for the container
// registry.
type Store interface {
	PutObject(key string, data io.Reader) error
	GetObject(key string) (io.ReadCloser, error)
	DeleteObject(key string) error
	ListObjects(key string) ([]string, error)
	FindObject(key string) (bool, error)
}
