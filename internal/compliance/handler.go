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

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	var body checkRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	departureDate, err := time.Parse("2006-01-02", body.DepartureDate)
	if err != nil {
		http.Error(w, "invalid departure_date", http.StatusBadRequest)
		return
	}

	returnDate, err := time.Parse("2006-01-02", body.ReturnDate)
	if err != nil {
		http.Error(w, "invalid return_date", http.StatusBadRequest)
		return
	}

	passportExpiry, err := time.Parse("2006-01-02", body.PassportExpiry)
	if err != nil {
		http.Error(w, "invalid passport_expiry", http.StatusBadRequest)
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
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	result := h.service.Check(r.Context(), req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(result)
}
