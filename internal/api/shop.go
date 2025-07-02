package api

import (
	"flomart/internal/catalog/shop/handler"
	"flomart/internal/catalog/shop/service"
	"flomart/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterShopRoutes(r chi.Router, h *handler.Handler, s service.Service) {
	r.Post("/", h.CreateShop)
	r.Get("/", h.ListShop)
	r.Get("/{id}", h.GetShopByID)
	r.With(middleware.RequireShopOwnershipOrAdmin(s)).Put("/{id}", h.UpdateShop)
	r.With(middleware.RequireShopOwnershipOrAdmin(s)).Delete("/{id}", h.DeleteShop)
}
