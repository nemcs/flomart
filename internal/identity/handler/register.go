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
	defer r.Body.Close()

	var input identity.RegisterInput
	//json.Decode + валидация данных полученных
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Fatal(err)
		return
	}

	//Регистрация пользователя -> передаем в service
	id, err := h.service.RegisterUser(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	// Успешный ответ пользаку
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(id); err != nil {
		log.Fatal(err)
		return
	}
}
