package minio

type Minio struct {
}

func Init() Minio {
	m := Minio{}
	return m
}

func (m *Minio) PutObject() error {
	return nil
}

func (m *Minio) GetObject() error {
	return nil
}
