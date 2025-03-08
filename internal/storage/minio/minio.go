package minio

import "io"

type Minio struct {
}

func Init() *Minio {
	m := Minio{}
	return &m
}

func (m *Minio) PutObject(key string, data io.Reader) error {
	return nil
}

func (m *Minio) GetObject(key string) (io.ReadCloser, error) {
	return nil, nil
}

func (m *Minio) DeleteObject(key string) error {
	return nil
}

func (m *Minio) ListObjects(key string) ([]string, error) {
	return nil, nil
}

func (m *Minio) FindObject(key string) (bool, error) {
	return false, nil
}
