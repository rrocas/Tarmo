package system

type VersionResponse struct {
	Server     string `json:"server"`
	APIVersion int    `json:"apiVersion"`
	MinClient  string `json:"minClient"`
}
