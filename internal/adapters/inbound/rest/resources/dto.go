package resources

type ResourceJSONResponseDTO struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        int     `json:"price"`
	BaseQuantity float64 `json:"base_quantity"`
	BaseUnit     string  `json:"base_unit"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type CreateResourceRequestDTO struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       int     `json:"price"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
}

type UpdateResourceRequestDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       int     `json:"price"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
}
