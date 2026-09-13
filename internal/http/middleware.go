package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setupMiddleware(r *chi.Mux) {
	r.Use(middleware.Logger)
}
