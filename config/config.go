package config

type Config struct {
	DBUrl     string
	JWTSecret string
}

// TODO load .env github.com/joho/godotenv или os.Getenv():
func New() *Config {
	return &Config{
		DBUrl:     "postgres://flower:flower@postgres:5432/flower_db?sslmode=disable",
		JWTSecret: "secrGJ83gGihkwGKWu3gVDGbe5jkoDG3gpdA",
	}
}
