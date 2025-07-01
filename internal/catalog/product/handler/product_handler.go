package handler

import (
	"flomart/domain/product"
	"flomart/internal/apperror"
	"flomart/internal/catalog/product/dto"
	"flomart/internal/catalog/product/service"
	"flomart/pkg/contextutils"
	"flomart/pkg/httphelper"
	"flomart/pkg/validation"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// TODO навести порядок в статусах
// TODO убрать зависимости из identity ( appErr)
type Handler struct {
	srv service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{srv: s}
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	input, err := httphelper.DecodeJSON[dto.ProductInputCreate](r)
	if err != nil {
		appErr := apperror.New(err, "Невалидный JSON", "json decode error", http.StatusBadRequest)
		httphelper.LogAndWriteError(w, *appErr)
		return
	}
	ci := *input

	claim, err := contextutils.GetUserFromCtx(ctx)
	if err != nil {
		appErr := apperror.New(err, "Не удалось создать товар", "service.CreateProduct: при получения shopID из контекста", http.StatusBadRequest)
		httphelper.LogAndWriteError(w, *appErr)
		return
	}

	id, err := h.srv.CreateProduct(r.Context(), string(claim.ShopID), ci)
	if err != nil {
		appErr := apperror.New(err, "Не удалось создать товар", "service.CreateProduct:", http.StatusBadRequest)
		httphelper.LogAndWriteError(w, *appErr)
		return
	}
	_ = httphelper.WriteJSONResponse(w, map[string]string{"productID": id}, http.StatusCreated)
}
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, URLParamProductID)
	if err := h.srv.DeleteProduct(r.Context(), id); err != nil {
		appErr := apperror.New(err, "Не удалось удалить товар", "service.DeleteProduct:", http.StatusBadRequest)
		httphelper.LogAndWriteError(w, *appErr)
		return
	}
	_ = httphelper.WriteJSONResponse(w, map[string]string{"status": "OK"}, http.StatusOK)
}

// TODO сделать возможность отправки клиентом не всех полей
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := chi.URLParam(r, URLParamProductID)

	input, err := httphelper.DecodeJSON[dto.ProductInput](r)
	if err != nil {
		appErr := apperror.New(err, "Невалидный JSON", "json decode error", http.StatusBadRequest)
		httphelper.LogAndWriteError(w, *appErr)
		return
	}
	pi := *input

	// TODO валидация вынести в функцию? Обрабатывать тоже внутри функции?
	// валидация input
	if err = validation.ValidateStruct(pi); err != nil {
		appErr := apperror.Wrap(err, "некорректный запрос", "validation failed", http.StatusBadRequest)
		httphelper.LogAndWriteError(w, *appErr)
		return
	}

	p, err := h.srv.UpdateProduct(r.Context(), id, pi)
	if err != nil {
		appErr := apperror.Wrap(err, "Не удалось обновить данные магазина", "service.UpdateShop: ошибка при обновлении данных магазина", http.StatusBadRequest)
		_ = httphelper.WriteJSONError(w, *appErr)
		return
	}

	_ = httphelper.WriteJSONResponse(w, map[string]product.Product{"data": p}, http.StatusOK)

}
func (h *Handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, URLParamProductID)
	p, err := h.srv.GetProductByID(r.Context(), id)
	if err != nil {
		appErr := apperror.New(err, "Не удалось получить продукт", "service.GetProductByID: ошибка при получении продукта", http.StatusBadRequest)
		_ = httphelper.WriteJSONError(w, *appErr)
		return
	}
	_ = httphelper.WriteJSONResponse(w, map[string]product.Product{"data": p}, http.StatusOK)
}
func (h *Handler) ListProductByShopID(w http.ResponseWriter, r *http.Request) {
	shopID := chi.URLParam(r, URLParamShopID)
	ps, err := h.srv.ListProductByShopID(r.Context(), shopID)
	if err != nil {
		appErr := apperror.Wrap(err, "Не удалось получить список продуктов", "service.ListProductByShopID: ошибка при получении продукта", http.StatusBadRequest)
		_ = httphelper.WriteJSONError(w, *appErr)
		return
	}
	_ = httphelper.WriteJSONResponse(w, map[string][]product.Product{"data": ps}, http.StatusOK)
}
