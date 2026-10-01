package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"tarmo"

	"tarmo/internal/adapters/inbound/rest/middleware"
	"tarmo/internal/adapters/inbound/rest/resources"
	"tarmo/internal/adapters/inbound/rest/system"
	"tarmo/internal/adapters/inbound/rest/templates"
	"tarmo/internal/adapters/outbound/persistence/sqlite"
	"tarmo/internal/config"
	resourceUseCase "tarmo/internal/core/resources/usecase"
	templateUseCase "tarmo/internal/core/templates/usecase"
	"tarmo/internal/lib/logger"
)

// @title Tarmo API
// @version 2.1.0
// @description Optimize and control processes based on templates.
// @host localhost:9136
// @BasePath /
func main() {
	cfg := config.Load()

	logger.Info("Starting Tarmo %s on port %s", tarmo.Version(), cfg.Port)

	// DB connection
	db, err := sqlite.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		logger.Fatal("failed to connect db: %v", err)
		return
	}

	// Repositories
	templateRepo := sqlite.NewTemplateRepository(db)
	resourceRepo := sqlite.NewResourceRepository(db)

	// Use cases
	templateUC := templateUseCase.NewTemplateUseCase(templateRepo)
	resUC := resourceUseCase.NewResourceUseCase(resourceRepo)

	// Handlers
	systemHandler := system.NewHandler(tarmo.Version())
	templateHandler := templates.NewHandler(templateUC)
	resourceHandler := resources.NewHandler(resUC)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logging)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/", func(r chi.Router) {
		system.RegisterRoutes(r, systemHandler)
		templates.RegisterRoutes(r, templateHandler)
		resources.RegisterRoutes(r, resourceHandler)
	})

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		logger.Fatal("server error: %v", err)
		return
	}
}
