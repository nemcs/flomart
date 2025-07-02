package middleware

import (
	"flomart/domain/user"
	"flomart/internal/apperror"
	"flomart/internal/catalog/shop/service"
	"flomart/pkg/contextutils"
	"flomart/pkg/httphelper"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaim, err := contextutils.GetUserFromCtx(r.Context())
		if err != nil || user.Role(userClaim.Role) != user.RoleAdmin {
			appErr := apperror.New(err, "недостаточно прав у пользователя", "пользователь не админ", http.StatusForbidden)
			httphelper.LogAndWriteError(w, *appErr)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireRoles(roles ...user.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userClaim, err := contextutils.GetUserFromCtx(r.Context())
			if err != nil {
				appErr := apperror.New(err, "недостаточно прав", "ошибка извлечения пользователя", http.StatusForbidden)
				httphelper.LogAndWriteError(w, *appErr)
				return
			}

			for _, role := range roles {
				if user.Role(userClaim.Role) == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			appErr := apperror.New(nil, "недостаточно прав", "роль не соответствует", http.StatusForbidden)
			httphelper.LogAndWriteError(w, *appErr)
		})
	}
}

func RequireShopOwnershipOrAdmin(checker service.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			shopID := chi.URLParam(r, "id")
			userClaim, err := contextutils.GetUserFromCtx(ctx)
			if err != nil {
				appErr := apperror.New(err, "недостаточно прав", "ошибка извлечения пользователя из контекста", http.StatusForbidden)
				httphelper.LogAndWriteError(w, *appErr)
				return
			}
			isOwner, err := checker.IsShopOwner(ctx, shopID, string(userClaim.UserID))
			if err != nil {
				appErr := apperror.New(err, "недостаточно прав", "проверке прав пользователя", http.StatusForbidden)
				httphelper.LogAndWriteError(w, *appErr)
			}

			if isOwner || user.Role(userClaim.Role) == user.RoleAdmin {
				next.ServeHTTP(w, r)
				return
			}

			appErr := apperror.New(nil, "недостаточно прав", "роль не соответствует", http.StatusForbidden)
			httphelper.LogAndWriteError(w, *appErr)

		})
	}
}
