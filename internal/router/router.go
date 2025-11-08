package router

import (
	"net/http"
	"todo-api/internal/database"
	"todo-api/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(db *database.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := handlers.New(db)

	r.Get("/health", h.Health)
	r.Get("/api/v1/hello", h.Hello)
	r.Get("/api/v1/users", h.ListUsers)
	r.Post("/api/v1/users", h.CreateUser)

	return r
}
