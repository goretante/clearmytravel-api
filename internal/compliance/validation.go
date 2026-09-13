package compliance

import (
	"errors"
	"strings"
)

var (
	ErrNationalityRequired    = errors.New("nationality_code is required")
	ErrDestinationRequired    = errors.New("destination_code is required")
	ErrDepartureRequired      = errors.New("departure_date is required")
	ErrReturnRequired         = errors.New("return_date is required")
	ErrPassportExpiryRequired = errors.New("passport_expiry is required")

	ErrInvalidNationality = errors.New("nationality_code must be 3 characters")
	ErrInvalidDestination = errors.New("destination_code must be 3 characters")
	ErrInvalidDateRange   = errors.New("return_date must be on or after departure_date")
)

func ValidateCheckRequest(req CheckRequest) error {
	if strings.TrimSpace(req.NationalityCode) == "" {
		return ErrNationalityRequired
	}

	if len(req.NationalityCode) != 3 {
		return ErrInvalidNationality
	}

	if strings.TrimSpace(req.DestinationCode) == "" {
		return ErrDestinationRequired
	}

	if len(req.DestinationCode) != 3 {
		return ErrInvalidDestination
	}

	if req.DepartureDate.IsZero() {
		return ErrDepartureRequired
	}

	if req.ReturnDate.IsZero() {
		return ErrReturnRequired
	}

	if req.PassportExpiry.IsZero() {
		return ErrPassportExpiryRequired
	}

	if req.ReturnDate.Before(req.DepartureDate) {
		return ErrInvalidDateRange
	}

	return nil
}
