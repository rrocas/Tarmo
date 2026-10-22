package usecase

import (
	"tarmo/internal/core/shared"
	"tarmo/internal/core/templates/domain"
	"tarmo/internal/core/templates/ports/inbound"
	"tarmo/internal/core/templates/ports/outbound"
)

type TemplateUseCase struct {
	repo outbound.TemplateRepositoryPort
}

func NewTemplateUseCase(repo outbound.TemplateRepositoryPort) *TemplateUseCase {
	return &TemplateUseCase{repo: repo}
}

func (uc *TemplateUseCase) GetAll() ([]inbound.TemplateDTO, error) {
	templates, err := uc.repo.FindAll()
	if err != nil {
		return []inbound.TemplateDTO{}, err
	}

	dtos := make([]inbound.TemplateDTO, 0, len(templates))
	for _, template := range templates {
		if template == nil {
			continue
		}

		qtyDTO := shared.QuantityDTO{
			Value: template.QuantityValue(),
			Unit: shared.UnitDTO{
				Name: template.QuantityUnitName(),
			},
		}

		dtos = append(dtos, inbound.TemplateDTO{
			ID:          template.ID(),
			Name:        template.Name(),
			Description: template.Description(),
			Quantity:    qtyDTO,
			Difficulty:  template.Difficulty(),
			UpdatedAt:   template.UpdatedAt(),
			CreatedAt:   template.CreatedAt(),
			Resources:   mapResourcesToDTO(template.Resources()),
		})
	}

	return dtos, nil
}

func (uc *TemplateUseCase) GetByID(id int) (inbound.TemplateDTO, error) {
	template, err := uc.repo.FindByID(id)
	if err != nil {
		return inbound.TemplateDTO{}, err
	}

	steps := make([]inbound.StepDTO, 0, len(template.Steps()))
	for _, step := range template.Steps() {
		steps = append(steps, inbound.StepDTO{
			Name:         step.Name(),
			Instructions: step.Instructions(),
			Order:        step.Order(),
		})
	}

	qtyDTO := shared.QuantityDTO{
		Value: template.QuantityValue(),
		Unit: shared.UnitDTO{
			Name: template.QuantityUnitName(),
		},
	}

	return inbound.TemplateDTO{
		ID:          template.ID(),
		Name:        template.Name(),
		Description: template.Description(),
		Quantity:    qtyDTO,
		Difficulty:  template.Difficulty(),
		UpdatedAt:   template.UpdatedAt(),
		CreatedAt:   template.CreatedAt(),
		Steps:       steps,
		Resources:   mapResourcesToDTO(template.Resources()),
	}, nil
}

func (uc *TemplateUseCase) Create(cmd inbound.CreateTemplateCommand) (int, error) {
	steps, err := mapStepsToDomain(cmd.Steps)
	if err != nil {
		return 0, err
	}

	resources, err := mapResourcesToDomain(cmd.Resources)
	if err != nil {
		return 0, err
	}

	template, err := domain.NewTemplate(cmd.Name, cmd.Quantity, cmd.Unit, cmd.Difficulty, steps, cmd.Description, resources)
	if err != nil {
		return 0, err
	}

	id, err := uc.repo.Save(template)
	return id, err
}

func (uc *TemplateUseCase) Update(cmd inbound.UpdateTemplateCommand) error {
	template, err := uc.repo.FindByID(cmd.ID)
	if err != nil {
		return err
	}

	steps, err := mapStepsToDomain(cmd.Steps)
	if err != nil {
		return err
	}

	resources, err := mapResourcesToDomain(cmd.Resources)
	if err != nil {
		return err
	}

	err = template.Update(cmd.Name, cmd.Quantity, cmd.Unit, cmd.Difficulty, steps, cmd.Description, resources)
	if err != nil {
		return err
	}

	return uc.repo.Update(template)
}

func (uc *TemplateUseCase) Delete(id int) error {
	return uc.repo.Remove(id)
}

func mapStepsToDomain(dtos []inbound.StepCommand) ([]domain.Step, error) {
	steps := make([]domain.Step, 0, len(dtos))
	for i, s := range dtos {
		step, err := domain.NewStep(s.Name, s.Instructions, i+1)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	return steps, nil
}

func mapResourcesToDomain(dtos []inbound.ResourceRefCommand) ([]domain.ResourceRef, error) {
	resources := make([]domain.ResourceRef, 0, len(dtos))
	for _, r := range dtos {
		ref, err := domain.NewResourceRef(r.ResourceID, r.Quantity, r.Unit)
		if err != nil {
			return nil, err
		}
		resources = append(resources, ref)
	}
	return resources, nil
}

func mapResourcesToDTO(resources []domain.ResourceRef) []inbound.ResourceRefDTO {
	dtos := make([]inbound.ResourceRefDTO, 0, len(resources))
	for _, r := range resources {

		qtyDTO := shared.QuantityDTO{
			Value: r.QuantityValue(),
			Unit: shared.UnitDTO{
				Name: r.QuantityUnitName(),
			},
		}

		dtos = append(dtos, inbound.ResourceRefDTO{
			ResourceID: r.ResourceID(),
			Quantity:   qtyDTO,
		})
	}
	return dtos
}
