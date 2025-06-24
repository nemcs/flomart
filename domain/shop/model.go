package shop

import "flomart/domain/user"

type ID int
type Shop struct {
	ID          ID      `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Location    string  `json:"location,omitempty"`
	OwnerId     user.ID `json:"owner_id,omitempty"`
}
