package compliance

import (
	"errors"
	"testing"
	"time"
)

func validCheckRequest() CheckRequest {
	return CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC),
		PassportExpiry:  time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestValidateCheckRequestValid(t *testing.T) {
	req := validCheckRequest()

	if err := ValidateCheckRequest(req); err != nil {
		t.Fatalf("expected valid request, got %v", err)
	}
}

func TestValidateCheckRequestMissingNationality(t *testing.T) {
	req := validCheckRequest()
	req.NationalityCode = ""

	err := ValidateCheckRequest(req)

	if !errors.Is(err, ErrNationalityRequired) {
		t.Fatalf(
			"expected ErrNationalityRequired, got %v",
			err,
		)
	}
}

func TestValidateCheckRequestInvalidNationality(t *testing.T) {
	req := validCheckRequest()
	req.NationalityCode = "HR"

	err := ValidateCheckRequest(req)

	if !errors.Is(err, ErrInvalidNationality) {
		t.Fatalf(
			"expected ErrInvalidNationality, got %v",
			err,
		)
	}
}

func TestValidateCheckRequestMissingDestination(t *testing.T) {
	req := validCheckRequest()
	req.DestinationCode = ""

	err := ValidateCheckRequest(req)

	if !errors.Is(err, ErrDestinationRequired) {
		t.Fatalf(
			"expected ErrDestinationRequired, got %v",
			err,
		)
	}
}

func TestValidateCheckRequestInvalidDestination(t *testing.T) {
	req := validCheckRequest()
	req.DestinationCode = "US"

	err := ValidateCheckRequest(req)

	if !errors.Is(err, ErrInvalidDestination) {
		t.Fatalf(
			"expected ErrInvalidDestination, got %v",
			err,
		)
	}
}

func TestValidateCheckRequestMissingDeparture(t *testing.T) {
	req := validCheckRequest()
	req.DepartureDate = time.Time{}

	err := ValidateCheckRequest(req)

	if !errors.Is(err, ErrDepartureRequired) {
		t.Fatalf(
			"expected ErrDepartureRequired, got %v",
			err,
		)
	}
}

func TestValidateCheckRequestMissingReturn(t *testing.T) {
	req := validCheckRequest()
	req.ReturnDate = time.Time{}

	err := ValidateCheckRequest(req)

	if !errors.Is(err, ErrReturnRequired) {
		t.Fatalf(
			"expected ErrReturnRequired, got %v",
			err,
		)
	}
}

func TestValidateCheckRequestMissingPassportExpiry(t *testing.T) {
	req := validCheckRequest()
	req.PassportExpiry = time.Time{}

	err := ValidateCheckRequest(req)

	if !errors.Is(err, ErrPassportExpiryRequired) {
		t.Fatalf(
			"expected ErrPassportExpiryRequired, got %v",
			err,
		)
	}
}

func TestValidateCheckRequestInvalidDateRange(t *testing.T) {
	req := validCheckRequest()

	req.DepartureDate = time.Date(
		2026, 10, 20, 0, 0, 0, 0, time.UTC,
	)

	req.ReturnDate = time.Date(
		2026, 10, 1, 0, 0, 0, 0, time.UTC,
	)

	err := ValidateCheckRequest(req)

	if !errors.Is(err, ErrInvalidDateRange) {
		t.Fatalf(
			"expected ErrInvalidDateRange, got %v",
			err,
		)
	}
}
