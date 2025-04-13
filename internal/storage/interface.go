package storage

import (
	"github.com/google/uuid"
	"io"
)

/*

1. Put Chunk (put single chunk)
2. Complete Layer (complete a layer of chunks)
3. Put Layer (put entire layer all at once)
4. Get Layer
5. List Layers
6. Find Layer
7. Delete Layer

*/

// Store is the interface for the object store for the container registry
type Store interface {
	PutChunk(offset string, length int, id uuid.UUID, data io.Reader) error
	CompleteLayer() error
	PutLayer(key string, data io.Reader) error
	GetLayer(key string) (io.ReadCloser, error)
	DeleteLayer(key string) error
	ListLayers(key string) ([]string, error)
	FindLayer(key string) bool
	CreateImage() error
	DeleteImage() error
	GetImage() error
	FindImage() error
}
