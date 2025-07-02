package config

import (
	"flomart/pkg/logger"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	Server  Server
	DB      DBManager
	JWT     TokenManager
	Storage StorageManager
}
type Server struct {
	ListenAddress string
}
type DBManager struct {
	Url string
}
type TokenManager struct {
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiryHour  int
	RefreshTokenExpiryHour int
}
type StorageManager struct {
	StorageAddress  string
	StoragePassword string
	StorageDB       string
}

func New() *Config {
	return &Config{
		Server: Server{
			ListenAddress: os.Getenv("LISTEN_ADDRESS"),
		},
		DB: DBManager{
			Url: os.Getenv("DB_URL"),
		},
		JWT: TokenManager{
			AccessTokenSecret:      os.Getenv("ACCESS_TOKEN_SECRET"),
			RefreshTokenSecret:     os.Getenv("REFRESH_TOKEN_SECRET"),
			AccessTokenExpiryHour:  mustGetInt("ACCESS_TOKEN_EXPIRY_HOUR"),
			RefreshTokenExpiryHour: mustGetInt("REFRESH_TOKEN_EXPIRY_HOUR"),
		},
		Storage: StorageManager{
			StorageAddress:  os.Getenv("STORAGE_ADDRESS"),
			StoragePassword: os.Getenv("STORAGE_PASSWORD"),
			StorageDB:       os.Getenv("STORAGE_DB"),
		},
	}
}

func mustGetInt(key string) int {
	valueStr := os.Getenv(key)
	hours, err := strconv.Atoi(valueStr)
	if err != nil || hours <= 0 {
		logger.Log.Error("Ошибка жизни токена", slog.String(logger.FieldErr, err.Error()))
		os.Exit(1)
	}
	return hours
}
