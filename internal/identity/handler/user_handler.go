// HTTP-хендлеры
package handler

import (
	"flomart/domain/user"
	"flomart/internal/apperror"
	"flomart/internal/identity"
	"flomart/internal/identity/dto"
	"flomart/internal/identity/service"
	"flomart/internal/identity/utils"
	"flomart/pkg/httphelper"
	"net/http"
)

// TODO ⚠️	Нет return после _ = utils.WriteError в branch’ах — иногда забыл.
// TODO 🔧	utils.LogAndWriteError логирует всегда Warn. Делай уровень динамическим (Error для 5xx).
// TODO 💡	Валидацию можно вынести в middleware‑валидатор, чтобы убрать дубли в каждом хендлере.
type Handler struct {
	service service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{service: s}
}

//TODO надо ли везде LogAndWriteError чтобы было едино все или можно местами и просто WriteError, если логи не нужны

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	input, err := httphelper.DecodeJSON[dto.RegisterInput](r)
	if err != nil {
		appErr := apperror.New(err, identity.ErrInvalidJSONMsg, identity.ErrInvalidJSONDev, http.StatusBadRequest)
		utils.LogAndWriteError(w, *appErr)
		return
	}
	regInput := *input

	// валидация input
	if err = identity.ValidateRegisterInput(regInput); err != nil {
		appErr := apperror.Wrap(err, identity.ErrBadRequestMsg, identity.ErrValidationDev, http.StatusBadRequest)
		utils.LogAndWriteError(w, *appErr)
		return
	}

	//Регистрация пользователя -> передаем в service
	id, err := h.service.RegisterUser(ctx, regInput)
	if err != nil {
		appErr := apperror.Wrap(err, "Регистрация не удалась", "service.RegisterUser: ошибка регистрации", http.StatusBadRequest)
		_ = utils.WriteJSONError(w, *appErr)
		return
	}

	// Успешный ответ пользаку
	_ = utils.WriteJSONResponse(w, map[string]user.ID{"id": id}, http.StatusCreated)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	input, err := httphelper.DecodeJSON[dto.LoginInput](r)
	if err != nil {
		appErr := apperror.New(err, identity.ErrInvalidJSONMsg, identity.ErrInvalidJSONDev, http.StatusBadRequest)
		utils.LogAndWriteError(w, *appErr)
		return
	}

	token, err := h.service.LoginUser(ctx, *input)
	if err != nil {
		appErr := apperror.Wrap(err, identity.ErrInvalidCredentialsMsg, "service.LoginUser: неверный логин или пароль", http.StatusUnauthorized)
		_ = utils.WriteJSONError(w, *appErr)
		return
	}
	_ = utils.WriteJSONResponse(w, map[string]string{"token": token}, http.StatusOK)
}
