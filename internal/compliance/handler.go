package compliance

import (
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type checkRequest struct {
	NationalityCode string `json:"nationality_code"`
	DestinationCode string `json:"destination_code"`
	DepartureDate   string `json:"departure_date"`
	ReturnDate      string `json:"return_date"`
	PassportExpiry  string `json:"passport_expiry"`
	VisaStatus      string `json:"visa_status"`
	ETAStatus       string `json:"eta_status"`
}

func writeBadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	var body checkRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	// Required fields must be checked before parsing dates.
	if body.NationalityCode == "" {
		writeBadRequest(w, "nationality_code is required")
		return
	}

	if body.DestinationCode == "" {
		writeBadRequest(w, "destination_code is required")
		return
	}

	if body.DepartureDate == "" {
		writeBadRequest(w, "departure_date is required")
		return
	}

	if body.ReturnDate == "" {
		writeBadRequest(w, "return_date is required")
		return
	}

	if body.PassportExpiry == "" {
		writeBadRequest(w, "passport_expiry is required")
		return
	}

	departureDate, err := time.Parse(
		"2006-01-02",
		body.DepartureDate,
	)
	if err != nil {
		writeBadRequest(w, "invalid departure_date")
		return
	}

	returnDate, err := time.Parse(
		"2006-01-02",
		body.ReturnDate,
	)
	if err != nil {
		writeBadRequest(w, "invalid return_date")
		return
	}

	passportExpiry, err := time.Parse(
		"2006-01-02",
		body.PassportExpiry,
	)
	if err != nil {
		writeBadRequest(w, "invalid passport_expiry")
		return
	}

	req := CheckRequest{
		NationalityCode: body.NationalityCode,
		DestinationCode: body.DestinationCode,
		DepartureDate:   departureDate,
		ReturnDate:      returnDate,
		PassportExpiry:  passportExpiry,
		VisaStatus:      VisaStatus(body.VisaStatus),
		ETAStatus:       ETAStatus(body.ETAStatus),
	}

	if err := ValidateCheckRequest(req); err != nil {
		writeBadRequest(w, err.Error())
		return
	}

	result := h.service.Check(r.Context(), req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(result)
}
