package system

import (
	"encoding/json"
	"net/http"
)

const (
	apiVersion = 2
	minClient  = "2.0.0"
)

type Handler struct {
	version string
}

func NewHandler(version string) *Handler {
	return &Handler{version: version}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(VersionResponse{
		Server:     h.version,
		APIVersion: apiVersion,
		MinClient:  minClient,
	})
}
