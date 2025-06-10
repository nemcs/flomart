// service - Бизнес-логика
// Регистрация пользователя
package service

import (
	"context"
	"flomart/config"
	"flomart/domain/user"
	"flomart/internal/identity"
	"flomart/internal/identity/repository"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

type Service interface {
	RegisterUser(ctx context.Context, input identity.RegisterInput) (user.ID, error)
	LoginUser(ctx context.Context, input identity.LoginInput) (string, error)
}
type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

// TODO хещирование пароля с помощью bcrypt (чтобы харнить хэш, а не сам пароль это безопаснее)
// hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
// Валидируем данные, передаем в repository слой и принимаем от него id, который возвращаем в handler
func (s *service) RegisterUser(ctx context.Context, input identity.RegisterInput) (user.ID, error) {
	if !strings.Contains(input.Email, "@") {
		return "", fmt.Errorf("Некорректный email")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	u := user.New(input.Email, string(password), input.Role, input.Name, input.Phone)

	id, err := s.repo.CreateUser(ctx, *u)
	if err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}
	return id, nil
}

func (s *service) LoginUser(ctx context.Context, input identity.LoginInput) (string, error) {
	//находим пользователя
	u, err := s.repo.FindByEmail(ctx, input.Email)
	if u.Email != input.Email {
		return "", fmt.Errorf("No username found")
	}
	if err != nil {
		return "", err
	}

	//проверяем пароль
	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)); err != nil {
		return "", fmt.Errorf("Неверный пароль: %v", err)
	}

	// генерируем токен
	tokenStr, err := CreateToken()
	if err != nil {
		return "", fmt.Errorf("ошибка генерации JWT токена: %v", err)
	}
	return tokenStr, nil
}

func CreateToken() (string, error) {
	cfg := config.New()
	token := jwt.New(jwt.SigningMethodHS256)
	tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

//func VarifyToken(tokenStr string) error {
//	token, err := jwt.Parse(tokenStr)
//
//	return nil
//}
