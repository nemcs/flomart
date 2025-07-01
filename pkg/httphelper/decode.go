package httphelper

import (
	"encoding/json"
	"flomart/internal/identity"
	"fmt"
	"net/http"
)

//TODO ⚠️	В error‑wrap передаёшь identity.ErrInvalidJSON — хорошо,
//но в хендлере ты потерял errors.Is, используешь строку.
//Нужен errors.Is(err, identity.ErrInvalidJSON) для маршрутизации.

func DecodeJSON[T any](r *http.Request) (*T, error) {
	var t T
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		//TODO вынести из identity
		return nil, fmt.Errorf("%w: %v", identity.ErrInvalidJSON, err)
	}
	return &t, nil
}
