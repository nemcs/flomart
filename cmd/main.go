package main

import (
	"flomart/config"
	"flomart/internal/api"
	"flomart/internal/identity"
	"flomart/pkg/db"
	"flomart/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
)

// TODO os.Exit или что-то иное? Или документировать код ошибки 1, 2 и т.д.
// TODO разобраться какие ошибки отдавать пользователю (безопасность) + когда warn/error
// TODO при рефреше старый access токен убивать
// TODO добавить больше инфы в логи при ошибках id юзера, магазина и прочая поебень

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Log.Error(identity.ErrLoadingEnvDev, slog.String(logger.FieldErr, err.Error()))
		os.Exit(1)
	}

	cfg := config.New()

	pool := db.NewPgxPool(cfg.DB.Url)
	defer pool.Close()

	identityHnd, shopHnd, shopSrv, productHnd := api.SetupDependencies(pool)

	r := chi.NewRouter()
	api.RegisterRoutes(r, cfg, identityHnd, shopHnd, shopSrv, productHnd)

	if err = http.ListenAndServe(cfg.Server.ListenAddress, r); err != nil {
		logger.Log.Error(identity.ErrRunServerDev, slog.String(logger.FieldErr, err.Error()))
		//пользователя никак не уведомляем об этом?
		os.Exit(1)
	}

}
