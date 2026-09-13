package trips

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/goretante/clearmytravel-api/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type tripRequest struct {
	UserID          string `json:"user_id"`
	PassportID      string `json:"passport_id"`
	DestinationCode string `json:"destination_code"`
	DepartureDate   string `json:"departure_date"`
	ReturnDate      string `json:"return_date"`
}

type tripResponse struct {
	ID              string `json:"id"`
	UserID          string `json:"user_id"`
	PassportID      string `json:"passport_id"`
	DestinationCode string `json:"destination_code"`
	DepartureDate   string `json:"departure_date"`
	ReturnDate      string `json:"return_date"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body tripRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	userID, err := parseUUID(body.UserID)
	if err != nil {
		writeBadRequest(w, "invalid user_id")
		return
	}

	passportID, err := parseUUID(body.PassportID)
	if err != nil {
		writeBadRequest(w, "invalid passport_id")
		return
	}

	departureDate, err := parseDate(body.DepartureDate)
	if err != nil {
		writeBadRequest(w, "invalid departure_date")
		return
	}

	returnDate, err := parseDate(body.ReturnDate)
	if err != nil {
		writeBadRequest(w, "invalid return_date")
		return
	}

	if body.DestinationCode == "" {
		writeBadRequest(w, "destination_code is required")
		return
	}

	params := db.CreateTripParams{
		UserID:          userID,
		PassportID:      passportID,
		DestinationCode: body.DestinationCode,
		DepartureDate:   departureDate,
		ReturnDate:      returnDate,
	}

	trip, err := h.service.Create(r.Context(), params)
	if err != nil {
		if errors.Is(err, ErrInvalidDateRange) {
			writeBadRequest(w, err.Error())
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, toTripResponse(trip))
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeBadRequest(w, "invalid trip id")
		return
	}

	trip, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "trip not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, toTripResponse(trip))
}

func (h *Handler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUUID(chi.URLParam(r, "userID"))
	if err != nil {
		writeBadRequest(w, "invalid user id")
		return
	}

	trips, err := h.service.ListByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]tripResponse, 0, len(trips))

	for _, trip := range trips {
		response = append(response, toTripResponse(trip))
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeBadRequest(w, "invalid trip id")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseUUID(value string) (pgtype.UUID, error) {
	var id pgtype.UUID

	if err := id.Scan(value); err != nil {
		return pgtype.UUID{}, err
	}

	return id, nil
}

func parseDate(value string) (pgtype.Date, error) {
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return pgtype.Date{}, err
	}

	return pgtype.Date{
		Time:  date,
		Valid: true,
	}, nil
}

func toTripResponse(trip db.Trip) tripResponse {
	return tripResponse{
		ID:              uuidString(trip.ID),
		UserID:          uuidString(trip.UserID),
		PassportID:      uuidString(trip.PassportID),
		DestinationCode: trip.DestinationCode,
		DepartureDate:   trip.DepartureDate.Time.Format("2006-01-02"),
		ReturnDate:      trip.ReturnDate.Time.Format("2006-01-02"),
	}
}

func uuidString(value pgtype.UUID) string {
	if !value.Valid {
		return ""
	}

	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		value.Bytes[0:4],
		value.Bytes[4:6],
		value.Bytes[6:8],
		value.Bytes[8:10],
		value.Bytes[10:16],
	)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func writeBadRequest(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusBadRequest, map[string]string{
		"error": message,
	})
}
