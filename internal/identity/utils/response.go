package utils

import (
	"encoding/json"
	"flomart/pkg/logger"
	"log/slog"
	"net/http"
)

type APIResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type APIError struct {
	Error string `json:"error"`
}

type DevError struct {
	Msg string
	Err error
}

func WriteJSONError(w http.ResponseWriter, msg string, code int) error {
	return WriteJSONResponse(w, APIError{Error: msg}, code)
}

func WriteJSONResponse(w http.ResponseWriter, data any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func LogAndWriteError(w http.ResponseWriter, handlerName, userMsg string, devErr DevError, code int) {
	logger.Log.Warn(devErr.Msg,
		slog.String(logger.FieldErr, devErr.Err.Error()),
		slog.String(logger.Handler, handlerName),
	)
	_ = WriteJSONError(w, userMsg, code)
}
