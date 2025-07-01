package main

import (
	"flomart/config"
	catalogProductHandler "flomart/internal/catalog/product/handler"
	catalogProductRepository "flomart/internal/catalog/product/repository"
	catalogProductService "flomart/internal/catalog/product/service"
	catalogShopHandler "flomart/internal/catalog/shop/handler"
	catalogShopRepository "flomart/internal/catalog/shop/repository"
	catalogShopService "flomart/internal/catalog/shop/service"
	"flomart/internal/identity"
	identityHandler "flomart/internal/identity/handler"
	identityRepository "flomart/internal/identity/repository"
	identityService "flomart/internal/identity/service"
	"flomart/internal/middleware"
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
	//TODO вынести в конфиг или куда?
	err := godotenv.Load(".env")
	if err != nil {
		logger.Log.Error(identity.ErrLoadingEnvDev, slog.String(logger.FieldErr, err.Error()))
		os.Exit(1)
	}

	cfg := config.New()

	pool := db.NewPgxPool(cfg.DBUrl)
	defer pool.Close()

	// === Identity ===
	identityRepo := identityRepository.NewRepository(pool)
	identitySrv := identityService.NewService(identityRepo)
	identityHnd := identityHandler.NewHandler(identitySrv)

	// === Catalog → Shop ===
	shopRepo := catalogShopRepository.NewRepository(pool)
	shopSrv := catalogShopService.NewService(shopRepo, pool)
	shopHnd := catalogShopHandler.NewHandler(shopSrv)

	// === Catalog → Product ===
	productRepo := catalogProductRepository.NewRepository(pool)
	productSrv := catalogProductService.NewService(productRepo)
	productHnd := catalogProductHandler.NewHandler(productSrv)

	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", identityHnd.RegisterUser)
		r.Post("/login", identityHnd.LoginUser)
		r.Post("/refresh", identityHnd.RefreshTokens)

	})
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.AccessTokenSecret))
		r.Get("/profile", identityHnd.ProfileUser)
		r.Route("/shops", func(r chi.Router) {
			r.Post("/", shopHnd.CreateShop)
			r.Get("/", shopHnd.ListShop)
			r.Get("/{id}", shopHnd.GetShopByID)
			r.With(middleware.RequireShopOwnershipOrAdmin(shopSrv)).Put("/{id}", shopHnd.UpdateShop)
			r.With(middleware.RequireShopOwnershipOrAdmin(shopSrv)).Delete("/{id}", shopHnd.DeleteShop)

		})
		r.Route("/products", func(r chi.Router) {
			r.Post("/", productHnd.CreateProduct)                             //CreateProduct
			r.Get("/{id}", productHnd.GetProductByID)                         //GetProductByID
			r.Get("/shops/{shopID}/products", productHnd.ListProductByShopID) //ListProductByShopID
			r.Put("/{id}", productHnd.UpdateProduct)                          //UpdateProduct
			r.Delete("/{id}", productHnd.DeleteProduct)                       // +++                   //DeleteProduct
		})
	})

	if err = http.ListenAndServe(cfg.ListenAddress, r); err != nil {
		logger.Log.Error(identity.ErrRunServerDev, slog.String(logger.FieldErr, err.Error()))
		//пользователя никак не уведомляем об этом?
		os.Exit(1)
	}

}
