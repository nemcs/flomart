// HTTP-хендлеры
package handler

import (
	"encoding/json"
	"flomart/internal/identity"
	"flomart/internal/identity/service"
	"log"
	"net/http"
)

type Handler struct {
	service service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{service: s}
}

// TODO обработку ошибок
// Обработка POST /register
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var input identity.RegisterInput
	//json.Decode + валидация данных полученных
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	//Регистрация пользователя -> передаем в service
	id, err := h.service.RegisterUser(ctx, input)
	if err != nil {
		//TODO возвращать json {"error": "invalid payload"}
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error: %v", err)
		return
	}

	// Успешный ответ пользаку
	// TODO возвращать json {"id": "uuid"}
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(id); err != nil {
		log.Printf("Error: %v", err)
		return
	}
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	//в сервис
	ctx := r.Context()
	defer r.Body.Close()

	//Принимаем email + password и парсим
	var input identity.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	// вызов сервиса
	token, err := h.service.LoginUser(ctx, input)
	if err != nil {
		//TODO возвращать json {"error": "invalid payload"}
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error: %v", err)
		return
	}

	// проверка статус код 200
	// возврат jwt токена
	// TODO {"token": "fwegwegwgwegweg"}
	if err = json.NewEncoder(w).Encode(token); err != nil {
		http.Error(w, "Ошибка при кодировании json ", http.StatusInternalServerError)
	}
}
