package dto

// Ответы клиенту: IDResponse, TokenResponse и т.п.

type IDResponse struct {
	ID string `json:"id"`
}
type TokenResponse struct {
	Token string `json:"token"`
}
type APIResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
