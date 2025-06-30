package handler

import (
	"flomart/domain/shop"
	"flomart/internal/apperror"
	"flomart/internal/catalog/shop/dto"
	"flomart/internal/catalog/shop/service"
	"flomart/internal/identity"
	"flomart/pkg/httphelper"
	"flomart/pkg/validation"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	srv service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{srv: s}
}

func (h Handler) CreateShop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	appErr := requireAdmin(ctx)
	if appErr != nil {
		httphelper.LogAndWriteError(w, *appErr)
		return
	}

	input, err := httphelper.DecodeJSON[dto.CreateInput](r)
	if err != nil {
		appErr = apperror.New(err, identity.ErrInvalidJSONMsg, identity.ErrInvalidJSONDev, http.StatusBadRequest)
		httphelper.LogAndWriteError(w, *appErr)
		return
	}
	createInput := *input

	// TODO валидация
	// валидация input
	if err = validation.ValidateStruct(createInput); err != nil {
		appErr = apperror.Wrap(err, "некорректный запрос", "validation failed", http.StatusBadRequest)
		httphelper.LogAndWriteError(w, *appErr)
		return
	}
	//передаем в сервис
	id, err := h.srv.CreateShop(ctx, createInput)
	if err != nil {
		appErr = apperror.Wrap(err, "Не удалось создать магазин", "service.CreateShop: ошибка при создании магазина", http.StatusBadRequest)
		_ = httphelper.WriteJSONError(w, *appErr)
		return
	}

	_ = httphelper.WriteJSONResponse(w, map[string]shop.ID{"id": id}, http.StatusCreated)
}

func (h Handler) GetShopByID(w http.ResponseWriter, r *http.Request) {
	shopID := chi.URLParam(r, "id")
	fmt.Println("\n\n\nid: %v\n\n\n", shopID)
	shp, err := h.srv.GetShopByID(r.Context(), shop.ID(shopID))
	if err != nil {
		appErr := apperror.Wrap(err, "Не удалось получить магазин", "service.GetShopByID: ошибка при получении магазина", http.StatusBadRequest)
		_ = httphelper.WriteJSONError(w, *appErr)
		return
	}

	_ = httphelper.WriteJSONResponse(w, map[string]shop.Shop{"data": shp}, http.StatusOK)
}

func (h Handler) ListShop(w http.ResponseWriter, r *http.Request) {
	listShop, err := h.srv.ListShop(r.Context())
	if err != nil {
		appErr := apperror.Wrap(err, "не удалось получить список магазинов", "service.ListShop: ошибка при получении списка магазинов", http.StatusInternalServerError)
		_ = httphelper.WriteJSONError(w, *appErr)
		return
	}
	_ = httphelper.WriteJSONResponse(w, map[string][]shop.Shop{"data": listShop}, http.StatusOK)

}

// TODO частичное обновление информации (только город или название магазина)
func (h Handler) UpdateShop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appErr := requireAdmin(ctx)
	if appErr != nil {
		httphelper.LogAndWriteError(w, *appErr)
		return
	}
	//валидирую ли?
	shopID := chi.URLParam(r, "id")

	input, err := httphelper.DecodeJSON[dto.UpdateInput](r)
	if err != nil {
		appErr = apperror.New(err, identity.ErrInvalidJSONMsg, identity.ErrInvalidJSONDev, http.StatusBadRequest)
		httphelper.LogAndWriteError(w, *appErr)
		return
	}
	updateInput := *input

	shp, err := h.srv.UpdateShop(ctx, shop.ID(shopID), updateInput)
	if err != nil {
		appErr = apperror.Wrap(err, "Не удалось обновить данные магазина", "service.UpdateShop: ошибка при обновлении данных магазина", http.StatusBadRequest)
		_ = httphelper.WriteJSONError(w, *appErr)
		return
	}

	_ = httphelper.WriteJSONResponse(w, map[string]shop.Shop{"data": shp}, http.StatusOK)
}

func (h Handler) DeleteShop(w http.ResponseWriter, r *http.Request) {
	appErr := requireAdmin(r.Context())
	if appErr != nil {
		httphelper.LogAndWriteError(w, *appErr)
		return
	}

	id := chi.URLParam(r, "id")
	if err := h.srv.DeleteShop(r.Context(), shop.ID(id)); err != nil {
		appErr = apperror.Wrap(err, "не удалось удалить магазин", "service.DeleteShop: ошибка при удалении магазина", http.StatusInternalServerError)
		_ = httphelper.WriteJSONError(w, *appErr)
	}
	_ = httphelper.WriteJSONResponse(w, map[string]string{"status": "deleted"}, http.StatusOK)
}
