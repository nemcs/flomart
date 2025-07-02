package repository

import (
	"context"
	"flomart/domain/shop"
	"flomart/domain/user"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	txRepository
	FindShopByName(ctx context.Context, name string) (*shop.Shop, error)
	FindCityByName(ctx context.Context, name string) (*shop.City, error)
	GetShopByID(ctx context.Context, id shop.ID) (shop.Shop, error)
	GetOwnerIDByShopID(ctx context.Context, shopID shop.ID) (user.ID, error)
	ListShop(ctx context.Context) ([]shop.Shop, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (repo *repository) FindShopByName(ctx context.Context, name string) (*shop.Shop, error) {
	var s shop.Shop
	query := `select id, name, description, location_id, owner_id from shops where name = $1`

	row := repo.db.QueryRow(ctx, query, name)
	err := row.Scan(
		&s.ID,
		&s.Name,
		&s.Description,
		&s.Location.ID,
		&s.OwnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти магазин %q: %w", name, err)
	}
	return &s, nil
}

func (repo *repository) FindCityByName(ctx context.Context, name string) (*shop.City, error) {
	query := `select id, name from cities where name = $1`
	var c shop.City
	if err := repo.db.QueryRow(ctx, query, name).Scan(&c.ID, &c.Name); err != nil {
		return nil, fmt.Errorf("не удалось найти город %q: %w", name, err)
	}
	return &c, nil
}

func (repo *repository) GetShopByID(ctx context.Context, id shop.ID) (shop.Shop, error) {
	sql := `select 
    			s.id AS shop_id,
    			s.name AS shop_name,
    			s.description,
    			s.owner_id,
    			l.id AS location_id,
    			c.id AS city_id,
    			c.name AS city_name
			from shops s
			join location l ON s.location_id = l.id
			join cities c ON l.city_id = c.id
			where s.id = $1`
	var s shop.Shop
	if err := repo.db.QueryRow(ctx, sql, id).Scan(
		&s.ID,
		&s.Name,
		&s.Description,
		&s.OwnerID,
		&s.Location.ID,
		&s.Location.City.ID,
		&s.Location.City.Name,
	); err != nil {
		//Плохая практика так возвращать? Забыл
		return shop.Shop{}, err
	}

	return s, nil
}

func (repo *repository) ListShop(ctx context.Context) ([]shop.Shop, error) {
	query := `select 
    			s.id AS shop_id,
    			s.name AS shop_name,
    			s.description,
    			s.owner_id,
    			l.id AS location_id,
    			c.id AS city_id,
    			c.name AS city_name
			from shops s
			join location l ON s.location_id = l.id
			join cities c ON l.city_id = c.id`
	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return []shop.Shop{}, err
	}
	defer rows.Close()

	var shp []shop.Shop
	for rows.Next() {
		var s shop.Shop
		if err = rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.OwnerID,
			&s.Location.ID,
			&s.Location.City.ID,
			&s.Location.City.Name,
		); err != nil {
			return []shop.Shop{}, fmt.Errorf("при сканировании строки магазина: %w", err)
		}
		shp = append(shp, s)
	}

	return shp, nil
}

// Поидее лучше возвращать shop.OwnerID , но надо переименовывать shop
func (repo *repository) GetOwnerIDByShopID(ctx context.Context, shopID shop.ID) (user.ID, error) {
	sql := `select owner_id from shops where id = $1`
	var ownerID user.ID
	if err := repo.db.QueryRow(ctx, sql, shopID).Scan(&ownerID); err != nil {
		return "", fmt.Errorf("не удалось найти магазин %q: %w", shopID, err)
	}
	return ownerID, nil
}
