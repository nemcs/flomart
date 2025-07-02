package api

import (
	"flomart/config"
	catalogProductHandler "flomart/internal/catalog/product/handler"
	catalogShopHandler "flomart/internal/catalog/shop/handler"
	"flomart/internal/catalog/shop/service"
	identityHandler "flomart/internal/identity/handler"
	"flomart/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, cfg *config.Config, identityHnd *identityHandler.Handler, shopHnd *catalogShopHandler.Handler, shopSrv service.Service, productHnd *catalogProductHandler.Handler) {

	// public auth routes
	r.Route("/auth", func(r chi.Router) {
		RegisterAuthRoutes(r, identityHnd)
	})

	// protected routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWT.AccessTokenSecret))
		r.Get("/profile", identityHnd.ProfileUser)

		r.Route("/shops", func(r chi.Router) {
			RegisterShopRoutes(r, shopHnd, shopSrv)

		})
		r.Route("/products", func(r chi.Router) {
			RegisterProductRoutes(r, productHnd)
		})
	})
}
