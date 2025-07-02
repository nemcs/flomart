package product

import (
	"github.com/google/uuid"
	"time"
)

// TODO вынести в common
type ID string
type Product struct {
	ID          string    `json:"id"`
	ShopID      string    `json:"shop_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Available   bool      `json:"available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProduct(shopID string, name, description string, price float64, available bool) *Product {
	return &Product{
		ID:          uuid.New().String(),
		ShopID:      shopID,
		Name:        name,
		Description: description,
		Price:       price,
		Available:   available,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
