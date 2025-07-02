package api

import (
	catalogProductHandler "flomart/internal/catalog/product/handler"
	catalogProductRepository "flomart/internal/catalog/product/repository"
	catalogProductService "flomart/internal/catalog/product/service"
	catalogShopHandler "flomart/internal/catalog/shop/handler"
	catalogShopRepository "flomart/internal/catalog/shop/repository"
	catalogShopService "flomart/internal/catalog/shop/service"
	identityHandler "flomart/internal/identity/handler"
	identityRepository "flomart/internal/identity/repository"
	identityService "flomart/internal/identity/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDependencies(pool *pgxpool.Pool) (*identityHandler.Handler, *catalogShopHandler.Handler, catalogShopService.Service, *catalogProductHandler.Handler) {
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

	return identityHnd, shopHnd, shopSrv, productHnd
}
