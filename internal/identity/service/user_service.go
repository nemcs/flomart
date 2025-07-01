// service - Бизнес-логика
// Регистрация пользователя
package service

import (
	"context"
	"errors"
	"flomart/config"
	"flomart/domain/user"
	"flomart/internal/identity"
	"flomart/internal/identity/dto"
	"flomart/internal/identity/repository"
	"flomart/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

//TODO ⚠️	Ветка existing != nil → возвращаешь ErrEmailAlreadyExists, но клиента информируешь строкой из ошибки. Лучше вернуть AppError (код 409) сразу из сервиса.
//TODO ⚠️	В ветке repo.CreateUser — пропускаешь pgErr.Code == 23505. TODO уже есть — сделай.
//TODO 🔧	logger.Log.Warn для ошибки БД — всё‑таки это Error (лёгкая потеря данных).
//TODO 💡	Пароли: добавь bcrypt.MinCost в конфиг, чтобы можно было менять cost.

type Service interface {
	RegisterUser(ctx context.Context, input dto.RegisterInput) (user.ID, error)
	LoginUser(ctx context.Context, input dto.LoginInput) (dto.TAccessToken, dto.TRefreshToken, error)
	RefreshTokens(ctx context.Context, input dto.RefreshInput) (dto.TAccessToken, dto.TRefreshToken, error)
}
type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterUser(ctx context.Context, input dto.RegisterInput) (user.ID, error) {
	existing, err := s.repo.FindByEmail(ctx, input.Email)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.Log.Warn("Ошибка при проверке email", slog.String(logger.FieldErr, err.Error()))
		return "", identity.ErrEmailSearch
	}
	if existing != nil {
		logger.Log.Warn("Email уже зарегистрирован", slog.String("email", input.Email))
		return "", identity.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Warn("Ошибка при генерации хэша пароля", slog.String(logger.FieldErr, err.Error()))
		return "", identity.ErrHashingPassword
	}

	u := user.New(input.Email, string(hashedPassword), input.Role, input.Name, input.Phone)

	id, err := s.repo.CreateUser(ctx, *u)
	//TODO вынести эту проверку в repo
	//проверка на pgError.Code == "23505"
	// и возвращать ошибку что пользователь уже есть
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		return "", identity.ErrEmailAlreadyExists
	}
	if err != nil {
		logger.Log.Error("Пользователь не создан", slog.String(logger.FieldSQL, err.Error()))
		return "", err
	}
	return *id, nil
}

func (s *service) LoginUser(ctx context.Context, input dto.LoginInput) (dto.TAccessToken, dto.TRefreshToken, error) {
	//находим пользователя
	u, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		logger.Log.Info(identity.ErrUserNotFoundDev, slog.String(logger.FieldSQL, err.Error()))
		return "", "", identity.ErrUserNotFound
	}

	//проверяем пароль
	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)); err != nil {
		logger.Log.Info(identity.ErrPasswordHashDev, slog.String(logger.FieldErr, err.Error()))
		return "", "", identity.ErrInvalidCredentials
	}

	shopID, err := s.repo.GetShopIDByUserID(ctx, u.ID)
	if err != nil {
		return "", "", errors.New("ошибка получения айди магазина")
	}

	// генерируем токен
	access, refresh, err := identity.CreateTokens(u.ID, shopID, u.Role)
	if err != nil {
		logger.Log.Error(identity.ErrTokenGenDev, slog.String(logger.FieldErr, err.Error()))
		//TODO возвращать ErrInternalServerMsg
		return "", "", identity.ErrGeneratingJWT
	}

	return access, refresh, nil
}

// TODO когда будет Redis нужно будет просто чекать не истек ли refresh?
func (s *service) RefreshTokens(ctx context.Context, input dto.RefreshInput) (dto.TAccessToken, dto.TRefreshToken, error) {
	// TODO лучше прокинуть в сервис, дабы не вызывать каждый раз на лету
	cfg := config.New()

	// 1. Парсим токен
	claims, err := identity.ParseToken(string(input.RefreshToken), cfg.RefreshTokenSecret)
	if err != nil {
		return "", "", identity.ErrInvalidToken
	}

	// 2. Проверяем, существует ли пользователь
	u, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return "", "", identity.ErrUserNotFound
	}

	// 3. Генерируем новую пару токенов
	access, refresh, err := identity.CreateTokens(u.ID, claims.ShopID, u.Role)
	if err != nil {
		logger.Log.Error(identity.ErrTokenGenDev, slog.String(logger.FieldErr, err.Error()))
		//TODO возвращать ErrInternalServerMsg
		return "", "", identity.ErrGeneratingJWT
	}

	return access, refresh, nil
}
