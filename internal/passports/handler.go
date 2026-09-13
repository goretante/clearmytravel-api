package passports

import (
	"encoding/json"
	"net/http"
	"time"
	"uuid"

	"github.com/go-chi/chi/v5"
	"github.com/goretante/clearmytravel-api/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	queries *db.Queries
}

func NewHandler(queries *db.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

type createPassportRequest struct {
	UserID          string `json:"user_id"`
	NationalityCode string `json:"nationality_code"`
	ExpiresAt       string `json:"expires_at"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPassportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	expiresAt, err := time.Parse("2006-01-02", req.ExpiresAt)
	if err != nil {
		http.Error(w, "invalid expires_at", http.StatusBadRequest)
		return
	}

	passport, err := h.queries.CreatePassport(r.Context(), db.CreatePassportParams{
		UserID: pgtype.UUID{
			Bytes: userID,
			Valid: true,
		},
		NationalityCode: req.NationalityCode,
		ExpiresAt: pgtype.Date{
			Time:  expiresAt,
			Valid: true,
		},
	})
	if err != nil {
		http.Error(w, "failed to create passport", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(passport)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	passportID := chi.URLParam(r, "id")

	id, err := uuid.Parse(passportID)
	if err != nil {
		http.Error(w, "invalid passport id", http.StatusBadRequest)
		return
	}

	passport, err := h.queries.GetPassportByID(r.Context(), pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
	if err != nil {
		http.Error(w, "passport not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(passport)
}

func (h *Handler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	id, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	passports, err := h.queries.ListPassportsByUserID(r.Context(), pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
	if err != nil {
		http.Error(w, "failed to get passports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(passports)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	passportID := chi.URLParam(r, "id")

	id, err := uuid.Parse(passportID)
	if err != nil {
		http.Error(w, "invalid passport id", http.StatusBadRequest)
		return
	}

	var req createPassportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	expiresAt, err := time.Parse("2006-01-02", req.ExpiresAt)
	if err != nil {
		http.Error(w, "invalid expires_at", http.StatusBadRequest)
		return
	}

	passport, err := h.queries.UpdatePassport(r.Context(), db.UpdatePassportParams{
		ID: pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		NationalityCode: req.NationalityCode,
		ExpiresAt: pgtype.Date{
			Time:  expiresAt,
			Valid: true,
		},
	})
	if err != nil {
		http.Error(w, "failed to update passport", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(passport)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	passportID := chi.URLParam(r, "id")

	id, err := uuid.Parse(passportID)
	if err != nil {
		http.Error(w, "invalid passport id", http.StatusBadRequest)
		return
	}

	err = h.queries.DeletePassport(r.Context(), pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
	if err != nil {
		http.Error(w, "failed to delete passport", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
