package service

import (
	"context"
	"errors"
	"flomart/domain/shop"
	"flomart/internal/catalog/shop/dto"
	"flomart/internal/catalog/shop/repository"
	"flomart/pkg/logger"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Service interface {
	CreateShop(ctx context.Context, input dto.CreateInput) (shop.ID, error)
	ListShop(ctx context.Context) ([]shop.Shop, error)
	GetShopByID(ctx context.Context, id shop.ID) (shop.Shop, error)
	DeleteShop(ctx context.Context, id shop.ID) error
	UpdateShop(ctx context.Context, shopID shop.ID, input dto.UpdateInput) (shop.Shop, error)
	IsShopOwner(ctx context.Context, shopID, userID string) (bool, error)
}

type service struct {
	repo repository.Repository
	db   *pgxpool.Pool
}

func NewService(repo repository.Repository, db *pgxpool.Pool) Service {
	return &service{repo: repo, db: db}
}

// вообще не нравится функция эта стремная
func (s *service) CreateShop(ctx context.Context, input dto.CreateInput) (shop.ID, error) {
	//Чекаем есть ли магазин в бд с таким же именем
	existing, err := s.repo.FindShopByName(ctx, input.Name)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.Log.Warn("ошибка при проверке имени магазина", slog.String(logger.FieldErr, err.Error()))
		return "", errors.New("имя магазина уже занято")
	}
	if existing != nil {
		logger.Log.Warn("Магазин уже существует", slog.String("name", input.Name))
		return "", errors.New("имя магазина уже занято")
	}

	//Чекаем есть ли город в нашей бд для смены города магазина
	city, err := s.repo.FindCityByName(ctx, input.CityName)
	if err != nil {
		logger.Log.Error("Город не найден", slog.String(logger.FieldErr, err.Error()))
		return "", errors.New("город не найден")
	}

	location := shop.NewLocation(*city)
	shp := shop.NewShop(input.Name, input.Description, shop.ID(input.OwnerID), *location)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err = s.repo.CreateLocationTx(ctx, tx, *location); err != nil {
		return "", err
	}
	id, err := s.repo.CreateShopTx(ctx, tx, *shp)
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		return "", errors.New("магазин уже существует / или имя занято?")
	}
	if err != nil {
		logger.Log.Error("Магазин не создан", slog.String(logger.FieldSQL, err.Error()))
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit tx: %w", err)
	}

	return *id, nil
}

func (s *service) ListShop(ctx context.Context) ([]shop.Shop, error) {
	shops, err := s.repo.ListShop(ctx)
	if err != nil {
		return []shop.Shop{}, fmt.Errorf("ошибка при получении списка магазинов: %w", err)
	}
	return shops, nil
}

func (s *service) GetShopByID(ctx context.Context, id shop.ID) (shop.Shop, error) {
	shp, err := s.repo.GetShopByID(ctx, id)
	if err != nil {
		return shop.Shop{}, fmt.Errorf("ошибка при получении магазина: %w", err)
	}
	return shp, nil
}

func (s *service) DeleteShop(ctx context.Context, id shop.ID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	locationID, err := s.repo.DeleteShopTx(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении магазина из таблицы Shops: %w", err)
	}
	err = s.repo.DeleteLocationTx(ctx, tx, locationID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении из таблицы Location %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (s *service) UpdateShop(ctx context.Context, shopID shop.ID, input dto.UpdateInput) (shop.Shop, error) {
	city, err := s.repo.FindCityByName(ctx, input.CityName)
	if err != nil {
		logger.Log.Error("Город не найден", slog.String(logger.FieldErr, err.Error()))
		return shop.Shop{}, errors.New("город не найден")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return shop.Shop{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err = s.repo.UpdateShopTx(ctx, tx, input.Name, input.Description, shopID, city.ID); err != nil {
		logger.Log.Error("не удалось обновить данные магазина", slog.String(logger.FieldErr, err.Error()))
		return shop.Shop{}, errors.New("не удалось обновить данные магазина")
	}

	if err = tx.Commit(ctx); err != nil {
		return shop.Shop{}, fmt.Errorf("commit tx: %w", err)
	}

	shp, err := s.repo.GetShopByID(ctx, shopID)
	if err != nil {
		logger.Log.Error("не удалось получить данные", slog.String(logger.FieldErr, err.Error()))
		return shop.Shop{}, errors.New("Данные обновлены, не удалось получить новые данные")
	}

	return shp, nil
}

func (s *service) IsShopOwner(ctx context.Context, shopID, userID string) (bool, error) {

	ownerID, err := s.repo.GetOwnerIDByShopID(ctx, shop.ID(shopID))
	if err != nil {
		logger.Log.Error("не удалось получить данные", slog.String(logger.FieldErr, err.Error()))
		return false, errors.New("ошибка при запросе в бд")
	}
	return string(ownerID) == userID, nil
}
