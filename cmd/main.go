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

// TODO сделать обработку ошибок
func main() {
	cfg := config.New()

	//TODO сделать отдельный startup CLI-процесс для миграций (например go run cmd/migrate/main.go)
	_, err := migrations.RunMigrations(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Ошибка при миграции: %v", err)
	}

	pool := db.NewPgxPool(cfg.DBUrl)
	defer pool.Close()

	repo := repository.NewRepository(pool)
	s := service.NewService(repo)
	h := handler.NewHandler(s)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/v1/register", h.RegisterUser)
	r.Post("/api/v1/login", h.LoginUser)

	log.Fatal(http.ListenAndServe(":8080", r))

}
