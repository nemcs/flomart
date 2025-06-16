package main

import (
	"flomart/config"
	"flomart/internal/identity"
	"flomart/internal/identity/handler"
	"flomart/internal/identity/repository"
	"flomart/internal/identity/service"
	"flomart/pkg/db"
	"flomart/pkg/logger"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
)

//TODO os.Exit или что-то иное? Или документировать код ошибки 1, 2 и т.д.
//TODO разобраться какие ошибки отдавать пользователю (безопасность) + когда warn/error
//TODO AppError ...............

func main() {
	//TODO вынести в конфиг или куда?
	err := godotenv.Load(".env")
	if err != nil {
		logger.Log.Error(identity.ErrLoadingEnvDev, slog.String(logger.FieldErr, err.Error()))
		os.Exit(1)
	}

	cfg := config.New()

	pool := db.NewPgxPool(cfg.DBUrl)
	defer pool.Close()

	repo := repository.NewRepository(pool)
	s := service.NewService(repo)
	h := handler.NewHandler(s)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.LoginUser)
	})

	port := fmt.Sprintf(":%s", cfg.Port)
	if err = http.ListenAndServe(port, r); err != nil {
		logger.Log.Error(identity.ErrRunServerDev, slog.String(logger.FieldErr, err.Error()))
		//пользователя никак не уведомляем об этом?
		os.Exit(1)
	}

}
