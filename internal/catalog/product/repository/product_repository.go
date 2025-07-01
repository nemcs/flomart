package repository

import (
	"context"
	"errors"
	"flomart/domain/product"
	"flomart/internal/catalog/product/dto"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrProductNotFound = errors.New("product not found")

type Repository interface {
	CreateProduct(ctx context.Context, p product.Product) (string, error) // productID, err
	GetProductByID(ctx context.Context, productID string) (product.Product, error)
	UpdateProduct(ctx context.Context, productID string, input dto.ProductInput) error
	DeleteProduct(ctx context.Context, productID string) error
	ListProductByShopID(ctx context.Context, shopID string) ([]product.Product, error)
}
type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (repo *repository) CreateProduct(ctx context.Context, p product.Product) (string, error) {
	sql := `insert into products(id, shop_id, name, description, price, available, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning id`
	var id string
	if err := repo.db.QueryRow(ctx, sql,
		p.ID, p.ShopID, p.Name, p.Description,
		p.Price, p.Available, p.CreatedAt, p.UpdatedAt).Scan(&id); err != nil {
		return "", fmt.Errorf("repo.CreateProduct Scan error: %w", err)
	}
	return id, nil
}
func (repo *repository) DeleteProduct(ctx context.Context, productID string) error {
	sql := `delete from products where id = $1`
	tag, err := repo.db.Exec(ctx, sql, productID)
	if err != nil {
		return fmt.Errorf("repo.DeleteProduct Exec error: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrProductNotFound
	}
	return nil
}
func (repo *repository) GetProductByID(ctx context.Context, productID string) (product.Product, error) {
	sql := `select id, shop_id, name, description, price, available, created_at, updated_at from products where id = $1`
	var p product.Product
	if err := repo.db.QueryRow(ctx, sql, productID).Scan(
		&p.ID,
		&p.ShopID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Available,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return product.Product{}, fmt.Errorf("repo.GetProductByID QueryRow error: %w", err)
	}
	return p, nil
}
func (repo *repository) ListProductByShopID(ctx context.Context, shopID string) ([]product.Product, error) {
	sql := `select id, shop_id, name, description, price, available, created_at, updated_at from products where shop_id = $1`
	rows, err := repo.db.Query(ctx, sql, shopID)
	if err != nil {
		return []product.Product{}, fmt.Errorf("repo.ListProductByShopID Query error: %w", err)
	}
	defer rows.Close()

	var ps []product.Product
	var p product.Product
	for rows.Next() {
		if err = rows.Scan(
			&p.ID,
			&p.ShopID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Available,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return []product.Product{}, fmt.Errorf("repo.ListProductByShopID Scan error: %w", err)
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (repo *repository) UpdateProduct(ctx context.Context, productID string, input dto.ProductInput) error {
	sql := `update products set name = $1, description = $2, price = $3, available = $4 where id = $5`
	tag, err := repo.db.Exec(ctx, sql, input.Name, input.Description, input.Price, input.Available, productID)
	if err != nil {
		return fmt.Errorf("repo.UpdateProduct Exec error: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrProductNotFound
	}
	return nil
}
