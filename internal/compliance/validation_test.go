package compliance

import (
	"testing"
	"time"
)

func validCheckRequest() CheckRequest {
	return CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2031, 6, 15, 0, 0, 0, 0, time.UTC),
		PassportExpiry:  time.Date(2032, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestValidateCheckRequestValid(t *testing.T) {
	req := validCheckRequest()

	if err := ValidateCheckRequest(req); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateCheckRequestMissingNationality(t *testing.T) {
	req := validCheckRequest()
	req.NationalityCode = ""

	if err := ValidateCheckRequest(req); err != ErrNationalityRequired {
		t.Fatalf("expected %v, got %v", ErrNationalityRequired, err)
	}
}

func TestValidateCheckRequestInvalidNationality(t *testing.T) {
	req := validCheckRequest()
	req.NationalityCode = "HR"

	if err := ValidateCheckRequest(req); err != ErrInvalidNationality {
		t.Fatalf("expected %v, got %v", ErrInvalidNationality, err)
	}
}

func TestValidateCheckRequestMissingDestination(t *testing.T) {
	req := validCheckRequest()
	req.DestinationCode = ""

	if err := ValidateCheckRequest(req); err != ErrDestinationRequired {
		t.Fatalf("expected %v, got %v", ErrDestinationRequired, err)
	}
}

func TestValidateCheckRequestInvalidDestination(t *testing.T) {
	req := validCheckRequest()
	req.DestinationCode = "US"

	if err := ValidateCheckRequest(req); err != ErrInvalidDestination {
		t.Fatalf("expected %v, got %v", ErrInvalidDestination, err)
	}
}

func TestValidateCheckRequestMissingDeparture(t *testing.T) {
	req := validCheckRequest()
	req.DepartureDate = time.Time{}

	if err := ValidateCheckRequest(req); err != ErrDepartureRequired {
		t.Fatalf("expected %v, got %v", ErrDepartureRequired, err)
	}
}

func TestValidateCheckRequestMissingReturn(t *testing.T) {
	req := validCheckRequest()
	req.ReturnDate = time.Time{}

	if err := ValidateCheckRequest(req); err != ErrReturnRequired {
		t.Fatalf("expected %v, got %v", ErrReturnRequired, err)
	}
}

func TestValidateCheckRequestMissingPassportExpiry(t *testing.T) {
	req := validCheckRequest()
	req.PassportExpiry = time.Time{}

	if err := ValidateCheckRequest(req); err != ErrPassportExpiryRequired {
		t.Fatalf("expected %v, got %v", ErrPassportExpiryRequired, err)
	}
}

func TestValidateCheckRequestInvalidDateRange(t *testing.T) {
	req := validCheckRequest()

	req.DepartureDate = time.Date(
		2031, 7, 1, 0, 0, 0, 0, time.UTC,
	)

	req.ReturnDate = time.Date(
		2031, 6, 15, 0, 0, 0, 0, time.UTC,
	)

	if err := ValidateCheckRequest(req); err != ErrInvalidDateRange {
		t.Fatalf("expected %v, got %v", ErrInvalidDateRange, err)
	}
}
