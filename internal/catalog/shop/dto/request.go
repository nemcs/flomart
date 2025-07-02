package dto

type ID string
type CreateInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	CityName    string `json:"city_name" validate:"required"`
	OwnerID     ID     `json:"owner_id" validate:"required"`
}
type IDInput struct {
	ID ID `json:"id"`
}

type UpdateInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	CityName    string `json:"city_name" validate:"required"`
}
