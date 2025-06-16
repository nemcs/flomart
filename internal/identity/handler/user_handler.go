// HTTP-хендлеры
package handler

import (
	"flomart/domain/user"
	"flomart/internal/identity"
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

	input, err := httphelper.DecodeJSON[identity.RegisterInput](r)
	if err != nil {
		utils.LogAndWriteError(w, "RegisterUser",
			identity.ErrInvalidJSONMsg,
			utils.DevError{
				Msg: identity.ErrInvalidJSONDev,
				Err: err,
			},
			http.StatusBadRequest)
		return
	}
	regInput := *input

	// валидация input
	if err = identity.ValidateRegisterInput(regInput); err != nil {
		utils.LogAndWriteError(w, "RegisterUser",
			identity.ErrBadRequestMsg,
			utils.DevError{
				Msg: identity.ErrValidationDev,
				Err: err,
			},
			http.StatusBadRequest)
		return
	}

	//Регистрация пользователя -> передаем в service
	id, err := h.service.RegisterUser(ctx, regInput)
	if err != nil {
		_ = utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Успешный ответ пользаку
	_ = utils.WriteJSONResponse(w, map[string]user.ID{"id": id}, http.StatusCreated)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	input, err := httphelper.DecodeJSON[identity.LoginInput](r)
	if err != nil {
		utils.LogAndWriteError(w, "LoginUser",
			identity.ErrInvalidJSONMsg,
			utils.DevError{
				Msg: identity.ErrInvalidJSONDev,
				Err: err,
			},
			http.StatusBadRequest)
		return
	}

	token, err := h.service.LoginUser(ctx, *input)
	if err != nil {
		_ = utils.WriteJSONError(w, identity.ErrInvalidCredentialsMsg, http.StatusUnauthorized)
		return
	}
	_ = utils.WriteJSONResponse(w, map[string]string{"token": token}, http.StatusOK)
}
