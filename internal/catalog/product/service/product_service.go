package service

import (
	"context"
	"flomart/domain/product"
	"flomart/internal/catalog/product/dto"
	"flomart/internal/catalog/product/repository"
	"fmt"
)

type Service interface {
	CreateProduct(ctx context.Context, shopID string, input dto.ProductInputCreate) (string, error) // productID, err
	GetProductByID(ctx context.Context, productID string) (product.Product, error)
	UpdateProduct(ctx context.Context, productID string, input dto.ProductInput) (product.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	ListProductByShopID(ctx context.Context, shopID string) ([]product.Product, error)
}
type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateProduct(ctx context.Context, shopID string, input dto.ProductInputCreate) (string, error) {

	p := product.NewProduct(shopID, input.Name, input.Description, input.Price, input.Available)

	id, err := s.repo.CreateProduct(ctx, *p)
	if err != nil {
		return "", fmt.Errorf("service.CreateProduct failed: %w", err)
	}
	return id, nil
}

func (s *service) GetProductByID(ctx context.Context, productID string) (product.Product, error) {
	p, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return product.Product{}, fmt.Errorf("service.GetProductByID failed: %w", err)
	}
	return p, nil
}

func (s *service) UpdateProduct(ctx context.Context, productID string, input dto.ProductInput) (product.Product, error) {
	if err := s.repo.UpdateProduct(ctx, productID, input); err != nil {
		return product.Product{}, fmt.Errorf("service.UpdateProduct failed: %w", err)
	}

	p, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return product.Product{}, fmt.Errorf("service.UpdateProduct failed in GetProductByID: %w", err)
	}

	return p, nil
}

func (s *service) DeleteProduct(ctx context.Context, productID string) error {
	if err := s.repo.DeleteProduct(ctx, productID); err != nil {
		return fmt.Errorf("service.DeleteProduct failed: %w", err)
	}
	return nil
}

func (s *service) ListProductByShopID(ctx context.Context, shopID string) ([]product.Product, error) {
	p, err := s.repo.ListProductByShopID(ctx, shopID)
	if err != nil {
		return []product.Product{}, fmt.Errorf("service.ListProductByShopID failed: %w", err)
	}
	return p, nil
}
