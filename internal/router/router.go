package router

import (
	"log"
	"net/http"
	"todo-api/internal/config"
	"todo-api/internal/handlers"
	auth "todo-api/internal/middleware"

	"todo-api/internal/repository/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(db *postgres.DB, cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := handlers.New(db, cfg)

	r.Get("/health", h.Health)
	r.Get("/api/v1/hello", h.Hello)
	r.Post("/api/v1/login", h.Login)
	// Users routes
	r.Group(func(r chi.Router) {
		r.Use(auth.Authenticate(cfg.JWT.Secret))
		r.Get("/api/v1/users", h.ListUsers)
		r.Post("/api/v1/users", h.CreateUser)
		r.Get("/api/v1/me", h.GetProfile)
		r.Get("/api/v1/users/{id}", h.FindUserByID)
		log.Printf("registered route: GET /api/v1/users/{id}")
		r.Get("/api/v1/users/{email}", h.FindUserByEmail)
		log.Printf("registered route: GET /api/v1/users/{email}")
	})

	// Tasks routes
	r.Route("/api/v1/tasks", func(r chi.Router) {
		r.Post("/", h.CreateTask)
		r.Get("/", h.ListTasks)
		r.Get("/{id}", h.GetTaskByID)
		r.Put("/{id}", h.UpdateTask)
		r.Delete("/{id}", h.DeleteTask)
	})

	return r
}
