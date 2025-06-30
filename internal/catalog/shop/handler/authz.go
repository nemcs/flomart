package handler

import (
	"context"
	"flomart/domain/user"
	"flomart/internal/apperror"
	"flomart/pkg/contextutils"
	"net/http"
)

func requireAdmin(ctx context.Context) *apperror.AppError {
	userClaim, err := contextutils.GetUserFromCtx(ctx)
	if err != nil {
		return apperror.New(err, "внутренняя ошибка", "ошибка при обращении к контексту", http.StatusInternalServerError)
	}
	if user.Role(userClaim.Role) != user.RoleAdmin {
		return apperror.New(nil, "недостаточно прав", "у клиента недостаточно прав", http.StatusForbidden)
	}
	return nil
}
