package utils

import (
	"encoding/json"
	"flomart/internal/apperror"
	"flomart/internal/identity/dto"
	"flomart/pkg/logger"
	"log/slog"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, data any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func WriteJSONError(w http.ResponseWriter, appErr apperror.AppError) error {
	return WriteJSONResponse(w, dto.APIError{Error: appErr.UserMsg}, appErr.Code)
}

func LogAndWriteError(w http.ResponseWriter, appErr apperror.AppError) {
	logger.Log.Warn(appErr.DevMsg, slog.String(logger.FieldErr, appErr.Err.Error()))
	_ = WriteJSONError(w, appErr)
}
