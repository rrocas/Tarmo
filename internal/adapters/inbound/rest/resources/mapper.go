package resources

import (
	"tarmo/internal/core/resources/ports/inbound"
	"time"
)

func ToCreateCommand(req CreateResourceRequestDTO) inbound.CreateResourceCommand {
	return inbound.CreateResourceCommand{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
	}
}

func ToUpdateCommand(req UpdateResourceRequestDTO) inbound.UpdateResourceCommand {
	return inbound.UpdateResourceCommand{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
	}
}

func ToResponse(dto *inbound.ResourceDTO) ResourceJSONResponseDTO {

	return ResourceJSONResponseDTO{
		ID:           dto.ID,
		Name:         dto.Name,
		Description:  dto.Description,
		Price:        dto.Price,
		BaseQuantity: dto.Quantity.Value,
		BaseUnit:     dto.Quantity.Unit.Name,
		CreatedAt:    dto.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    dto.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func ToResponseList(resources []inbound.ResourceDTO) []ResourceJSONResponseDTO {
	res := make([]ResourceJSONResponseDTO, 0, len(resources))
	for _, r := range resources {
		res = append(res, ToResponse(&r))
	}
	return res
}
