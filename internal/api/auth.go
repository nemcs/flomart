package api

import (
	"flomart/internal/identity/handler"
	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, h *handler.Handler) {
	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.LoginUser)
	r.Post("/refresh", h.RefreshTokens)
}
