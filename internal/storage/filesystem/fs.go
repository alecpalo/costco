package filesystem

import (
	"io"
	"os"
	"path/filepath"
)

type FileSystem struct {
	filePath string
}

func Init(fp string) *FileSystem {
	fs := FileSystem{
		filePath: fp,
	}
	return &fs
}

// PutObject takes in a key and a reader object and attempts to write the object
// to the provided location in the filesystem. Upon success returning nil, upon
// failure returning an error.
func (fs *FileSystem) PutObject(key string, data io.Reader) error {
	writePath := filepath.Join(fs.filePath, key)
	var bytes []byte

	_, err := data.Read(bytes)
	if err != nil {
		return err
	}

	err = os.WriteFile(writePath, bytes, os.ModePerm)

	return err
}

// GetObject takes in a key and attempts to read the file from that location
// returning a reader and an error.
func (fs *FileSystem) GetObject(key string) (io.ReadCloser, error) {
	readPath := filepath.Join(fs.filePath, key)

	return os.Open(readPath)
}

func (fs *FileSystem) DeleteObject(key string) error {
	return nil
}

func (fs *FileSystem) ListObjects(key string) ([]string, error) {
	return nil, nil
}

// FindObject attempts to find the file stored at key, returning true if
// is successfully found otherwise false or an error.
func (fs *FileSystem) FindObject(key string) bool {
	searchPath := filepath.Join(fs.filePath, key)

	_, err := os.Stat(searchPath)
	if err != nil {
		return false
	}

	return true
}
