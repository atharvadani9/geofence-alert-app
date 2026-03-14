package api

import (
	"net/http"

	"github.com/atharvadani9/geofence-alert-app/internal/api/handlers"
	"github.com/atharvadani9/geofence-alert-app/internal/api/middleware"
	"github.com/atharvadani9/geofence-alert-app/internal/config"
	"github.com/atharvadani9/geofence-alert-app/internal/database"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates and configures the HTTP router
func NewRouter(db *database.DB, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS)

	// Health check endpoint
	r.Get("/health", handlers.Health)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handlers.RegisterUser(db, cfg))
			r.Post("/login", handlers.LoginUser(db, cfg))
			// r.Post("/refresh", handlers.RefreshToken)
			// r.Get("/me", handlers.GetCurrentUser)
		})

		// Protected routes (to be implemented)
		// r.Group(func(r chi.Router) {
		// 	r.Use(middleware.Auth)
		// 	// Add protected routes here
		// })
	})

	return r
}
