package templates

type StepJSONResponseDTO struct {
	Order        int    `json:"order"`
	Name         string `json:"name"`
	Instructions string `json:"instructions"`
}

type ResourceRefResponseDTO struct {
	ResourceID int     `json:"resource_id"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit"`
}

type TemplateJSONResponseDTO struct {
	ID          int                      `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Quantity    float64                  `json:"quantity"`
	Unit        string                   `json:"unit"`
	Difficulty  int                      `json:"difficulty"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
	Steps       []StepJSONResponseDTO    `json:"steps"`
	Resources   []ResourceRefResponseDTO `json:"resources"`
}

type TemplateListJSONResponseDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateTemplateRequestDTO struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Quantity    float64          `json:"quantity"`
	Unit        string           `json:"unit"`
	Difficulty  int              `json:"difficulty"`
	Steps       []StepDTO        `json:"steps"`
	Resources   []ResourceRefDTO `json:"resources"`
}

type UpdateTemplateRequestDTO struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Quantity    float64          `json:"quantity"`
	Unit        string           `json:"unit"`
	Difficulty  int              `json:"difficulty"`
	Steps       []StepDTO        `json:"steps"`
	Resources   []ResourceRefDTO `json:"resources"`
}

type ResourceRefDTO struct {
	ResourceID int     `json:"resource_id"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit"`
}

type StepDTO struct {
	Name         string `json:"name"`
	Instructions string `json:"instructions"`
}
