package filesystem

import "io"

type FileSystem struct {
	filePath string
}

func Init(fp string) *FileSystem {
	fs := FileSystem{
		filePath: fp,
	}
	return &fs
}

func (fs *FileSystem) PutObject(key string, data io.Reader) error {
	return nil
}

func (fs *FileSystem) GetObject(key string) (io.ReadCloser, error) {
	return nil, nil
}

func (fs *FileSystem) DeleteObject(key string) error {
	return nil
}

func (fs *FileSystem) ListObjects(key string) ([]string, error) {
	return nil, nil
}

func (fs *FileSystem) FindObject(key string) (bool, error) {
	return false, nil
}
