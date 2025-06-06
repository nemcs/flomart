package config

type Config struct {
	DBUrl string
}

// TODO load .env
func New() *Config {
	return &Config{DBUrl: "postgres://flower:flower@postgres:5432/flower_db?sslmode=disable"}
}
