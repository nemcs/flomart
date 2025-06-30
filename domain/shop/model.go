package shop

import "github.com/google/uuid"

type ID string

type Shops struct {
	Shop []Shop `validate:"dive"` // Проверка каждого shop внутри слайса
}

type Shop struct {
	ID          ID       `json:"id,omitempty"`                             // ID генерируется, не валидируется
	Name        string   `json:"name" validate:"required,min=2,max=100"`   // Название обязательно, от 2 до 100 символов
	Description string   `json:"description" validate:"omitempty,max=500"` // Необязательное, макс. 500 символов
	OwnerID     ID       `json:"owner_id" validate:"required,uuid4"`       // ID владельца обязательно и должно быть UUIDv4
	Location    Location `json:"location" validate:"required,dive"`        // Обязательная структура Location
}

type Location struct {
	ID   ID   `json:"id" validate:"required,uuid4"`  // ID Location обязательно UUID
	City City `json:"city" validate:"required,dive"` // Обязательная вложенная структура
}

type City struct {
	ID   ID     `json:"id" validate:"required,uuid4"`   // UUID ID города
	Name string `json:"name" validate:"required,min=2"` // Обязательное имя, минимум 2 символа
}

//	func NewCity(name string) *City {
//		return &City{
//			ID:   ID(uuid.New().String()),
//			Name: name,
//		}
//	}
func NewLocation(city City) *Location {
	return &Location{
		ID:   ID(uuid.New().String()),
		City: city,
	}
}

func NewShop(name, description string, ownerID ID, location Location) *Shop {
	return &Shop{
		ID:          ID(uuid.New().String()),
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		Location: Location{
			ID: location.ID,
			City: City{
				ID:   location.City.ID,
				Name: location.City.Name,
			},
		},
	}
}
