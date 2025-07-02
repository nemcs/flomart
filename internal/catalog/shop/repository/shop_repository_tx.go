package repository

import (
	"context"
	"errors"
	"flomart/domain/shop"
	"flomart/internal/identity"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type txRepository interface {
	CreateShopTx(ctx context.Context, tx pgx.Tx, s shop.Shop) (*shop.ID, error)
	CreateLocationTx(ctx context.Context, tx pgx.Tx, loc shop.Location) error
	UpdateShopTx(ctx context.Context, tx pgx.Tx, name, description string, shopID, cityID shop.ID) error
	DeleteShopTx(ctx context.Context, tx pgx.Tx, shopID shop.ID) (shop.ID, error)
	DeleteLocationTx(ctx context.Context, tx pgx.Tx, locationID shop.ID) error
}

func (repo *repository) CreateLocationTx(ctx context.Context, tx pgx.Tx, loc shop.Location) error {
	//TODO переделать на Exec если не буду возвращать ничего
	sql := `insert into location (id, city_id) values ($1, $2)`
	tag, err := tx.Exec(ctx, sql, loc.ID, loc.City.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", identity.ErrSqlInsertDev, err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("location запись не создана")
	}
	return nil
}

func (repo *repository) CreateShopTx(ctx context.Context, tx pgx.Tx, s shop.Shop) (*shop.ID, error) {
	query := `insert into shops (name, description, location_id, owner_id) values ($1, $2, $3, $4) returning id`
	var id shop.ID
	if err := tx.QueryRow(ctx, query, s.Name, s.Description, s.Location.ID, s.OwnerID).Scan(&id); err != nil {
		return nil, fmt.Errorf("%s: %w", identity.ErrSqlInsertDev, err)
	}
	return &id, nil
}

func (repo *repository) UpdateShopTx(ctx context.Context, tx pgx.Tx, name, description string, shopID, cityID shop.ID) error {
	sqlUpNameDescr := `update shops set name = $1, description = $2 where id = $3 returning location_id`
	sqlUpCity := `update location set city_id = $1 where id = $2`

	var locationID shop.ID
	if err := tx.QueryRow(ctx, sqlUpNameDescr, name, description, shopID).Scan(&locationID); err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, sqlUpCity, cityID, locationID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("location таблица не обновлена")
	}

	return nil
}

func (repo *repository) DeleteShopTx(ctx context.Context, tx pgx.Tx, shopID shop.ID) (shop.ID, error) {
	sql := `delete from shops where id=$1 returning location_id`
	var locationID shop.ID
	if err := tx.QueryRow(ctx, sql, shopID).Scan(&locationID); err != nil {
		return "", fmt.Errorf("delete shop: %w", err)
	}
	return locationID, nil
}
func (repo *repository) DeleteLocationTx(ctx context.Context, tx pgx.Tx, locationID shop.ID) error {
	sql := `delete from location where id=$1`
	tag, err := tx.Exec(ctx, sql, locationID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("locationID не найден")
	}
	return nil
}
