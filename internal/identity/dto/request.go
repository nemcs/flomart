package dto

// Входящие структуры (например, LoginInput, RegisterInput — если не используются ещё где-то)

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone"`
	Role     string `json:"role" validate:"required,oneof=client courier admin"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
