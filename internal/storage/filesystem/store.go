package filesystem

import (
	"errors"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
)

type Store struct {
	filePath string
	fs       fs
}

func Init(fp string) *Store {
	fs := Store{
		filePath: fp,
		fs:       &filesystem{},
	}
	return &fs
}

// PutChunk puts a chunk into a staging location until all chunks have been created
// and expects the key, id that is uploaded to and the data itself. Upon success
// this function returns nil and upon failure it returns an error.
func (s *Store) PutChunk(offset string, length int, id uuid.UUID, data io.Reader) error {
	// todo: make this a real thing
	path := filepath.Join(s.filePath, "temp", id.String())
	info, err := s.fs.Stat(path)

	// todo make this a bit better
	if errors.Is(err, os.ErrNotExist) {
		err = s.fs.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if !info.IsDir() {
		return errors.New("file is not a directory")
	}

	offsetPath := filepath.Join(path, offset)
	f, err := s.fs.OpenFile(offsetPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	lenWritten, err := io.Copy(f, data)
	if err != nil {
		return err
	}

	if lenWritten != int64(length) {
		return errors.New("file length mismatch")
	}

	return nil
}

// PutObject takes in a key and a reader object and attempts to write the object
// to the provided location in the filesystem. Upon success returning nil, upon
// failure returning an error.
func (s *Store) PutObject(key string, data io.Reader) error {
	writePath := filepath.Join(s.filePath, key)
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
func (s *Store) GetObject(key string) (io.ReadCloser, error) {
	readPath := filepath.Join(s.filePath, key)

	return os.Open(readPath)
}

func (s *Store) DeleteObject(key string) error {
	return nil
}

func (s *Store) ListObjects(key string) ([]string, error) {
	return nil, nil
}

// FindObject attempts to find the file stored at key, returning true if
// is successfully found otherwise false or an error.
func (s *Store) FindObject(key string) bool {
	searchPath := filepath.Join(s.filePath, key)

	_, err := os.Stat(searchPath)
	if err != nil {
		return false
	}

	return true
}
