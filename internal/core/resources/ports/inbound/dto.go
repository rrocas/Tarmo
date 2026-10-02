package inbound

import (
	"tarmo/internal/core/shared"
	"time"
)

type CreateResourceCommand struct {
	Name        string
	Description string
	Price       int
	Quantity    float64
	Unit        string
}

type UpdateResourceCommand struct {
	ID          int
	Name        string
	Description string
	Price       int
	Quantity    float64
	Unit        string
}

type ResourceDTO struct {
	ID          int
	Name        string
	Description string
	Price       int
	Quantity    shared.QuantityDTO
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
