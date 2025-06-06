// service - Бизнес-логика
// Регистрация пользователя
package service

import (
	"flomart/domain/user"
	"flomart/internal/identity"
	"flomart/internal/identity/repository"
	"fmt"
	"log"
	"strings"
)

type Service interface {
	RegisterUser(input identity.RegisterInput) (user.ID, error)
}
type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

// Валидируем данные, передаем в repository слой и принимаем от него id, который возвращаем в handler
func (s *service) RegisterUser(input identity.RegisterInput) (user.ID, error) {
	if !strings.Contains(input.Email, "@") {
		return "", fmt.Errorf("Некорректный email")
	}
	id, err := s.repo.CreateUser(input)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return id, nil
}
