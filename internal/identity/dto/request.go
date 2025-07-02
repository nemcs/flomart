package dto

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`                     // Обязательный email в формате email
	Password string `json:"password" validate:"required,min=8"`                  // Обязательный, минимум 8 символов
	Name     string `json:"name" validate:"required"`                            // Обязательное имя
	Phone    string `json:"phone" validate:"omitempty,e164"`                     // Необязательный, но если указан — формат E.164 (например, +1234567890)
	Role     string `json:"role" validate:"required,oneof=client courier admin"` // Обязательный, одно из трёх значений
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"` // Обязательный email
	Password string `json:"password" validate:"required"`    // Обязательный пароль
}

type RefreshInput struct {
	RefreshToken TRefreshToken `json:"refresh_token" validate:"required"` // Обязательный refresh token
}
