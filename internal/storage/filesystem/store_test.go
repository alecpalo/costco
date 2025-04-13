package filesystem

import (
	"os"
	"testing"
)

type mockFileSystem struct {
	mockFile file
	fileErr  error
	fileInfo os.FileInfo
	statErr  error
	mkdirErr error
}

func (m *mockFileSystem) OpenFile(name string, flag int, perm os.FileMode) (file, error) {
	return os.OpenFile(name, flag, perm)
}

func (m *mockFileSystem) Open(name string) (file, error) {
	return os.Open(name)
}

func (m *mockFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (m *mockFileSystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// TestUploadChunk
func TestUploadChunk(t *testing.T) {

}

func TestUploadLayer(t *testing.T) {

}
