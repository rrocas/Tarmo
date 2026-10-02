package usecase

import (
	"tarmo/internal/core/resources/domain"
	"tarmo/internal/core/resources/ports/inbound"
	"tarmo/internal/core/resources/ports/outbound"
	"tarmo/internal/core/shared"
)

type ResourceUseCase struct {
	repo outbound.ResourceRepositoryPort
}

func NewResourceUseCase(repo outbound.ResourceRepositoryPort) *ResourceUseCase {
	return &ResourceUseCase{repo: repo}
}

func (uc *ResourceUseCase) GetAll() ([]inbound.ResourceDTO, error) {
	resources, err := uc.repo.FindAll()
	if err != nil {
		return []inbound.ResourceDTO{}, err
	}

	dtos := make([]inbound.ResourceDTO, 0, len(resources))
	for _, resource := range resources {
		if resource == nil {
			continue
		}

		qtyDTO := shared.QuantityDTO{
			Value: resource.QuantityValue(),
			Unit: shared.UnitDTO{
				Name: resource.QuantityUnitName(),
			},
		}

		dtos = append(dtos, inbound.ResourceDTO{
			ID:          resource.ID(),
			Name:        resource.Name(),
			Description: resource.Description(),
			Price:       resource.Price(),
			Quantity:    qtyDTO,
			CreatedAt:   resource.CreatedAt(),
			UpdatedAt:   resource.UpdatedAt(),
		})
	}

	return dtos, nil
}

func (uc *ResourceUseCase) GetByID(id int) (inbound.ResourceDTO, error) {
	resource, err := uc.repo.FindByID(id)
	if err != nil {
		return inbound.ResourceDTO{}, err
	}

	qtyDTO := shared.QuantityDTO{
		Value: resource.QuantityValue(),
		Unit: shared.UnitDTO{
			Name: resource.QuantityUnitName(),
		},
	}

	return inbound.ResourceDTO{
		ID:          resource.ID(),
		Name:        resource.Name(),
		Description: resource.Description(),
		Price:       resource.Price(),
		Quantity:    qtyDTO,
		CreatedAt:   resource.CreatedAt(),
		UpdatedAt:   resource.UpdatedAt(),
	}, nil
}

func (uc *ResourceUseCase) Create(cmd inbound.CreateResourceCommand) (int, error) {
	resource, err := domain.NewResource(cmd.Name, cmd.Description, cmd.Price, cmd.Quantity, cmd.Unit)
	if err != nil {
		return 0, err
	}

	id, err := uc.repo.Save(resource)
	return id, err
}

func (uc *ResourceUseCase) Update(cmd inbound.UpdateResourceCommand) error {
	resource, err := uc.repo.FindByID(cmd.ID)
	if err != nil {
		return err
	}

	err = resource.Update(cmd.Name, cmd.Description, cmd.Price, cmd.Quantity, cmd.Unit)
	if err != nil {
		return err
	}

	return uc.repo.Update(resource)
}

func (uc *ResourceUseCase) Delete(id int) error {
	return uc.repo.Remove(id)
}
