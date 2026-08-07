package db

type Config struct {
	// PostgresURI the db connection string for Postgres
	PostgresURI string
}

func NewConfig(connectionString string) *Config {
	return &Config{PostgresURI: connectionString}
}
