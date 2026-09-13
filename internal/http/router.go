package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goretante/clearmytravel-api/internal/compliance"
	"github.com/goretante/clearmytravel-api/internal/db"
	"github.com/goretante/clearmytravel-api/internal/passports"
	"github.com/goretante/clearmytravel-api/internal/trips"
	"github.com/goretante/clearmytravel-api/internal/users"
)

func NewRouter(queries *db.Queries, complianceService *compliance.Service) http.Handler {
	r := chi.NewRouter()

	setupMiddleware(r)

	userHandler := users.NewHandler(queries)
	passportHandler := passports.NewHandler(queries)
	complianceHandler := compliance.NewHandler(complianceService)
	tripHandler := trips.NewHandler(trips.NewService(trips.NewDBRepository(queries)))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// users routes
	r.Post("/users", userHandler.Create)
	r.Get("/users/{id}", userHandler.GetByID)
	r.Delete("/users/{id}", userHandler.Delete)

	// passport routes
	r.Post("/passports", passportHandler.Create)
	r.Get("/passports/{id}", passportHandler.GetByID)
	r.Get("/users/{userID}/passports", passportHandler.ListByUserID)
	r.Put("/passports/{id}", passportHandler.Update)
	r.Delete("/passports/{id}", passportHandler.Delete)

	// compliance routes
	r.Post("/compliance/check", complianceHandler.Check)

	// trips routes
	r.Post("/trips", tripHandler.Create)
	r.Get("/trips/{id}", tripHandler.GetByID)
	r.Get("/users/{userID}/trips", tripHandler.ListByUserID)
	r.Delete("/trips/{id}", tripHandler.Delete)

	return r
}
