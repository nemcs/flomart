package main

import (
	"flomart/config"
	"flomart/database/migrations"
	"flomart/internal/identity/handler"
	"flomart/internal/identity/repository"
	"flomart/internal/identity/service"
	"flomart/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	cfg := config.New()

	pool := db.NewPgxPool(cfg.DBUrl)
	defer pool.Close()

	repo := repository.NewRepository(pool)
	s := service.NewService(repo)
	h := handler.NewHandler(s)

	_, err := migrations.RunMigrations(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Ошибка при миграции: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/v1/register", h.RegisterUser)

	log.Fatal(http.ListenAndServe(":8080", r))

}
