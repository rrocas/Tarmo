package system

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Get("/ping", handler.Ping)
	r.Get("/version", handler.Version)
}
