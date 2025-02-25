package postgres

type PostgresMetadata struct {
}

func Init() PostgresMetadata {
	p := PostgresMetadata{}
	return p
}

func (m *PostgresMetadata) UpdateContainer() error {
	return nil
}
