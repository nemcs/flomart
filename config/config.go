package config

import (
	"flomart/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBUrl                  string
	ListenAddress          string
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiryHour  *jwt.NumericDate
	RefreshTokenExpiryHour *jwt.NumericDate
}

func New() *Config {
	return &Config{
		//TODO сделать проверку на пустую строку и fatal-лог
		DBUrl:         os.Getenv("DB_URL"),
		ListenAddress: os.Getenv("LISTEN_ADDRESS"),
		//TODO переместить в TokenManager
		AccessTokenSecret:      os.Getenv("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret:     os.Getenv("REFRESH_TOKEN_SECRET"),
		AccessTokenExpiryHour:  mustGetDurationHours("ACCESS_TOKEN_EXPIRY_HOUR"),
		RefreshTokenExpiryHour: mustGetDurationHours("REFRESH_TOKEN_EXPIRY_HOUR"),
	}
}

func mustGetDurationHours(key string) *jwt.NumericDate {
	valueStr := os.Getenv(key)
	hours, err := strconv.Atoi(valueStr)
	if err != nil || hours <= 0 {
		logger.Log.Error(logger.FieldErr, "Ошибка жизни токена")
		os.Exit(1)
	}
	return jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(hours)))

}
