package templates

import (
	"tarmo/internal/core/templates/ports/inbound"
	"time"
)

func ToCreateCommand(req CreateTemplateRequestDTO) inbound.CreateTemplateCommand {
	steps := make([]inbound.StepCommand, 0, len(req.Steps))
	for _, s := range req.Steps {
		steps = append(steps, inbound.StepCommand{
			Name:         s.Name,
			Instructions: s.Instructions,
		})
	}
	resources := make([]inbound.ResourceRefCommand, 0, len(req.Resources))
	for _, r := range req.Resources {
		resources = append(resources, inbound.ResourceRefCommand{
			ResourceID: r.ResourceID,
			Quantity:   r.Quantity,
			Unit:       r.Unit,
		})
	}
	return inbound.CreateTemplateCommand{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		Difficulty:  req.Difficulty,

		Steps:     steps,
		Resources: resources,
	}
}

func ToUpdateCommand(req UpdateTemplateRequestDTO) inbound.UpdateTemplateCommand {
	steps := make([]inbound.StepCommand, 0, len(req.Steps))
	for _, s := range req.Steps {
		steps = append(steps, inbound.StepCommand{
			Name:         s.Name,
			Instructions: s.Instructions,
		})
	}
	resources := make([]inbound.ResourceRefCommand, 0, len(req.Resources))
	for _, r := range req.Resources {
		resources = append(resources, inbound.ResourceRefCommand{
			ResourceID: r.ResourceID,
			Quantity:   r.Quantity,
			Unit:       r.Unit,
		})
	}
	return inbound.UpdateTemplateCommand{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		Difficulty:  req.Difficulty,
		Steps:       steps,
		Resources:   resources,
	}
}

func ToResponse(dto *inbound.TemplateDTO) TemplateJSONResponseDTO {
	steps := make([]StepJSONResponseDTO, len(dto.Steps))
	for i, s := range dto.Steps {
		steps[i] = StepJSONResponseDTO{
			Order:        s.Order,
			Name:         s.Name,
			Instructions: s.Instructions,
		}
	}

	resources := make([]ResourceRefResponseDTO, len(dto.Resources))
	for i, r := range dto.Resources {
		resources[i] = ResourceRefResponseDTO{
			ResourceID: r.ResourceID,
			Quantity:   r.Quantity.Value,
			Unit:       r.Quantity.Unit.Name,
		}
	}

	return TemplateJSONResponseDTO{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Quantity:    dto.Quantity.Value,
		Unit:        dto.Quantity.Unit.Name,
		Difficulty:  dto.Difficulty,
		CreatedAt:   dto.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   dto.UpdatedAt.UTC().Format(time.RFC3339),
		Steps:       steps,
		Resources:   resources,
	}
}

func ToListItemResponse(r inbound.TemplateDTO) TemplateListJSONResponseDTO {
	return TemplateListJSONResponseDTO{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
}

func ToResponseList(templates []inbound.TemplateDTO) []TemplateListJSONResponseDTO {
	res := make([]TemplateListJSONResponseDTO, 0, len(templates))
	for _, t := range templates {
		res = append(res, ToListItemResponse(t))
	}
	return res
}
