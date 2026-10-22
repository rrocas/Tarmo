package domain

import (
	"errors"
	"tarmo/internal/core/shared"
	"time"
)

// ===========================ERRORS============================

var (
	ErrNameRequired      = errors.New("name is required")
	ErrQuantityInvalid   = errors.New("quantity must be greater than 0")
	ErrUnitRequired      = errors.New("unit must be defined")
	ErrDifficultyInvalid = errors.New("difficulty must be between 0 and 5")
	ErrNoSteps           = errors.New("template must have at least one step")
	ErrInvalidStepOrders = errors.New("step orders must be sequential starting from 1")
	ErrStepName          = errors.New("step must have a name")
	ErrStepOrder         = errors.New("step order must be greater than 0")
	ErrResourceIDInvalid = errors.New("resource ID must be greater than 0")
)

// ===========================MODELS============================

type Template struct {
	id          int
	name        string // [REQUIRED]
	description string
	quantity    shared.Quantity // [REQUIRED]
	difficulty  int             // [REQUIRED]
	steps       []Step          // [REQUIRED] MIN 1
	resources   []ResourceRef
	createdAt   time.Time // [AUTO GENERATED] set by persistence
	updatedAt   time.Time // [AUTO GENERATED] set by persistence
}

type Step struct {
	order        int    // [REQUIRED]
	name         string // [REQUIRED]
	instructions string
}

type ResourceRef struct {
	resourceID int             // [REQUIRED]
	quantity   shared.Quantity // [REQUIRED]
}

// ===========================GETTERS===========================

func (t *Template) ID() int                   { return t.id }
func (t *Template) Name() string              { return t.name }
func (t *Template) Description() string       { return t.description }
func (t *Template) Quantity() shared.Quantity { return t.quantity }
func (t *Template) QuantityValue() float64    { return t.quantity.Value() }
func (t *Template) QuantityUnitName() string  { return t.quantity.Unit().Name }
func (t *Template) Difficulty() int           { return t.difficulty }
func (t *Template) Steps() []Step             { return append([]Step(nil), t.steps...) }
func (t *Template) Resources() []ResourceRef  { return append([]ResourceRef(nil), t.resources...) }
func (t *Template) CreatedAt() time.Time      { return t.createdAt }
func (t *Template) UpdatedAt() time.Time      { return t.updatedAt }

func (s *Step) Order() int           { return s.order }
func (s *Step) Name() string         { return s.name }
func (s *Step) Instructions() string { return s.instructions }

func (rr *ResourceRef) ResourceID() int          { return rr.resourceID }
func (rr *ResourceRef) QuantityValue() float64   { return rr.quantity.Value() }
func (rr *ResourceRef) QuantityUnitName() string { return rr.quantity.Unit().Name }

// ===========================CONSTRUCTORS======================

func NewResourceRef(resourceID int, quantity float64, unitStr string) (ResourceRef, error) {
	qty, err := shared.NewQuantity(quantity, unitStr)
	if err != nil {
		return ResourceRef{}, err
	}
	qty = qty.ToBase()
	rr := ResourceRef{
		resourceID: resourceID,
		quantity:   qty,
	}
	if err := rr.Validate(); err != nil {
		return ResourceRef{}, err
	}
	return rr, nil
}

func NewStep(name, instructions string, order int) (Step, error) {
	s := Step{
		name:         name,
		instructions: instructions,
		order:        order,
	}
	if err := s.Validate(); err != nil {
		return Step{}, err
	}
	return s, nil
}

func NewTemplate(
	name string,
	quantity float64,
	unitStr string,
	difficulty int,
	steps []Step,
	description string,
	resources []ResourceRef,
) (*Template, error) {

	qty, err := shared.NewQuantity(quantity, unitStr)
	if err != nil {
		return nil, err
	}
	qty = qty.ToBase()
	r := &Template{
		id:          0,
		name:        name,
		quantity:    qty,
		difficulty:  difficulty,
		steps:       steps,
		description: description,
		resources:   resources,
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	return r, nil
}

func ReconstructTemplate(
	id int,
	name string,
	quantity float64,
	unitStr string,
	difficulty int,
	steps []Step,
	description string,
	resources []ResourceRef,
	createdAt, updatedAt time.Time,
) (*Template, error) {

	qty, err := shared.NewQuantity(quantity, unitStr)
	if err != nil {
		return nil, err
	}
	qty = qty.ToBase()
	t := &Template{
		id:          id,
		name:        name,
		quantity:    qty,
		difficulty:  difficulty,
		steps:       steps,
		description: description,
		resources:   resources,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	return t, nil
}

// ===========================METHODS===========================
func (t *Template) Update(name string, quantity float64, unitStr string, difficulty int, steps []Step, description string, resources []ResourceRef) error {
	temp, err := NewTemplate(name, quantity, unitStr, difficulty, steps, description, resources)
	if err != nil {
		return err
	}

	t.name = temp.name
	t.quantity = temp.quantity
	t.difficulty = temp.difficulty
	t.steps = temp.steps
	t.description = temp.description
	t.resources = temp.resources

	return nil
}

func (t *Template) UpdateResources(resources []ResourceRef) error {
	for _, resource := range resources {
		if err := resource.Validate(); err != nil {
			return err
		}
	}
	t.resources = resources
	return nil
}

// ===========================VALIDATORS========================

func (t *Template) Validate() error {
	if t.name == "" {
		return ErrNameRequired
	}
	if t.difficulty < 0 || t.difficulty > 5 {
		return ErrDifficultyInvalid
	}
	if len(t.steps) == 0 {
		return ErrNoSteps
	}

	for i, step := range t.steps {
		expectedOrder := i + 1
		if step.order != expectedOrder {
			return ErrInvalidStepOrders
		}
		if err := step.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Step) Validate() error {

	if s.order <= 0 {
		return ErrStepOrder
	}

	if s.name == "" {
		return ErrStepName
	}

	return nil
}

func (rr *ResourceRef) Validate() error {
	if rr.resourceID <= 0 {
		return ErrResourceIDInvalid
	}
	return nil
}
