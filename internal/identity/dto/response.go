package dto

// Ответы клиенту: IDResponse, TokenResponse и т.п.

// TODO без T в go такое не принято
type TAccessToken string
type TRefreshToken string

type IDResponse struct {
	ID string `json:"id"`
}

//	type TokenResponse struct {
//		Token string `json:"token"`
//	}
type TokenPairResponse struct {
	AccessToken  TAccessToken  `json:"access_token"`
	RefreshToken TRefreshToken `json:"refresh_token"`
}
type APIResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
