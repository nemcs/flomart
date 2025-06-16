package config

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Config struct {
	Port          string
	DBUrl         string
	JWTSecret     string
	JWTExpiration *jwt.NumericDate
}

func New() *Config {
	return &Config{
		//TODO сделать проверку на пустую строку и fatal-лог
		Port:      os.Getenv("LISTEN_PORT"),
		DBUrl:     os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		//TODO переместить JWTExpiration в TokenManager
		JWTExpiration: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
	}
}
