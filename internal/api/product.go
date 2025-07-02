package api

import (
	"flomart/internal/catalog/product/handler"
	"github.com/go-chi/chi/v5"
)

func RegisterProductRoutes(r chi.Router, h *handler.Handler) {
	r.Post("/", h.CreateProduct)
	r.Get("/{id}", h.GetProductByID)
	r.Get("/shops/{shopID}/products", h.ListProductByShopID)
	r.Put("/{id}", h.UpdateProduct)
	r.Delete("/{id}", h.DeleteProduct)
}
