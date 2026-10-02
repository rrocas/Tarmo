package domain

import (
	"errors"
	"tarmo/internal/core/shared"
	"time"
)

// ===========================ERRORS============================

var (
	ErrNameRequired         = errors.New("name is required")
	ErrPriceRequired        = errors.New("price is required")
	ErrPriceGreaterThanZero = errors.New("price must be greater than 0")
	ErrUnitRequired         = errors.New("unit is required")
	ErrQuantityInvalid      = errors.New("quantity must be greater than 0")
)

// ===========================MODELS============================

type Resource struct {
	id          int             // [REQUIRED] [AUTO GENERATED]
	name        string          // [REQUIRED]
	description string          // [OPTIONAL]
	price       int             // [REQUIRED] In cents, or smallest currency unit
	quantity    shared.Quantity // [REQUIRED] e.g. 100, 1
	createdAt   time.Time       // [AUTO GENERATED] set by persistence
	updatedAt   time.Time       // [AUTO GENERATED] set by persistence
}

// ===========================GETTERS===========================

func (r *Resource) ID() int                   { return r.id }
func (r *Resource) Name() string              { return r.name }
func (r *Resource) Description() string       { return r.description }
func (r *Resource) Price() int                { return r.price }
func (r *Resource) Quantity() shared.Quantity { return r.quantity }
func (r *Resource) QuantityValue() float64    { return r.quantity.Value() }
func (r *Resource) QuantityUnitName() string  { return r.quantity.Unit().Name }
func (r *Resource) CreatedAt() time.Time      { return r.createdAt }
func (r *Resource) UpdatedAt() time.Time      { return r.updatedAt }

// ===========================CONSTRUCTORS======================

func NewResource(name string, description string, price int, quantity float64, unitStr string) (*Resource, error) {
	qty, err := shared.NewQuantity(quantity, unitStr)
	if err != nil {
		return nil, err
	}

	baseQty := qty.ToBase()

	rsc := &Resource{
		name:        name,
		description: description,
		price:       price,
		quantity:    baseQty,
	}
	if err := rsc.Validate(); err != nil {
		return nil, err
	}
	return rsc, nil
}

func ReconstructResource(id int, name string, description string, price int, quantity float64, unitStr string, createdAt, updatedAt time.Time) (*Resource, error) {
	qty, err := shared.NewQuantity(quantity, unitStr)
	if err != nil {
		return nil, err
	}

	baseQty := qty.ToBase()

	rsc := &Resource{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		quantity:    baseQty,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
	if err := rsc.Validate(); err != nil {
		return nil, err
	}
	return rsc, nil
}

// ===========================METHODS===========================

func (r *Resource) Update(name string, description string, price int, quantity float64, unitStr string) error {
	rsc, err := NewResource(name, description, price, quantity, unitStr)
	if err != nil {
		return err
	}

	r.name = rsc.name
	r.description = rsc.description
	r.price = rsc.price
	r.quantity = rsc.quantity

	return nil
}

// ===========================VALIDATORS========================

func (r *Resource) Validate() error {
	if r.name == "" {
		return ErrNameRequired
	}
	if r.price <= 0 {
		return ErrPriceGreaterThanZero
	}
	return nil
}
