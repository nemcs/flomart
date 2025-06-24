package main

import (
	"flomart/config"
	"flomart/internal/identity"
	"flomart/internal/identity/handler"
	"flomart/internal/identity/repository"
	"flomart/internal/identity/service"
	"flomart/internal/middleware"
	"flomart/pkg/db"
	"flomart/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
)

//TODO os.Exit или что-то иное? Или документировать код ошибки 1, 2 и т.д.
//TODO разобраться какие ошибки отдавать пользователю (безопасность) + когда warn/error

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
	//глобально или локально для защищенных роутеров? r.Use(middleware.AuthMiddleware(cfg.AccessTokenSecret))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.LoginUser)
		r.Post("/refresh", h.RefreshTokens)

	})
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.AccessTokenSecret))
		r.Get("/profile", h.ProfileUser)
		r.Route("/shops", func(r chi.Router) {
			r.Post("/", r.ServeHTTP)
			r.Get("/", r.ServeHTTP)
			r.Get("/:id", r.ServeHTTP)
			r.Put("/:id", r.ServeHTTP)
			r.Delete("/:id", r.ServeHTTP)

		})
	})

	if err = http.ListenAndServe(cfg.ListenAddress, r); err != nil {
		logger.Log.Error(identity.ErrRunServerDev, slog.String(logger.FieldErr, err.Error()))
		//пользователя никак не уведомляем об этом?
		os.Exit(1)
	}

}
