package main

import (
	"flomart/config"
	"flomart/database/migrations"
	"flomart/internal/identity"
	"flomart/pkg/logger"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func main() {
	//TODO вынести в конфиг или куда?
	err := godotenv.Load(".env")
	if err != nil {
		logger.Log.Error(identity.ErrLoadingEnvDev, slog.String(logger.FieldErr, err.Error()))
		os.Exit(1)
	}
	cfg := config.New()

	_, err = migrations.RunMigrations(cfg.DBUrl)
	if err != nil {
		logger.Log.Error(identity.ErrRunMigrationsDev, slog.String(logger.FieldErr, err.Error()))
		//пользователя никак не уведомляем об этом?
		os.Exit(1)
	}
}
