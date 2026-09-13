package trips

import (
	"context"
	"testing"
	"time"

	"github.com/goretante/clearmytravel-api/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type fakeRepository struct {
	createCalled bool
	trip         db.Trip
}

func (f *fakeRepository) Create(
	ctx context.Context,
	params db.CreateTripParams,
) (db.Trip, error) {
	f.createCalled = true

	return f.trip, nil
}

func (f *fakeRepository) GetByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.Trip, error) {
	return f.trip, nil
}

func (f *fakeRepository) ListByUserID(
	ctx context.Context,
	userID pgtype.UUID,
) ([]db.Trip, error) {
	return []db.Trip{f.trip}, nil
}

func (f *fakeRepository) ListByPassportID(
	ctx context.Context,
	passportID pgtype.UUID,
) ([]db.Trip, error) {
	return []db.Trip{f.trip}, nil
}

func (f *fakeRepository) Update(
	ctx context.Context,
	params db.UpdateTripParams,
) (db.Trip, error) {
	return f.trip, nil
}

func (f *fakeRepository) Delete(
	ctx context.Context,
	id pgtype.UUID,
) error {
	return nil
}

func TestServiceCreateRejectsInvalidDateRange(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)

	params := db.CreateTripParams{
		DepartureDate: pgtype.Date{
			Time:  time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		ReturnDate: pgtype.Date{
			Time:  time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	}

	_, err := service.Create(context.Background(), params)

	if err != ErrInvalidDateRange {
		t.Fatalf(
			"expected ErrInvalidDateRange, got %v",
			err,
		)
	}

	if repository.createCalled {
		t.Fatal("repository should not be called")
	}
}

func TestServiceCreateCallsRepository(t *testing.T) {
	repository := &fakeRepository{
		trip: db.Trip{},
	}

	service := NewService(repository)

	params := db.CreateTripParams{
		DepartureDate: pgtype.Date{
			Time:  time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		ReturnDate: pgtype.Date{
			Time:  time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	}

	_, err := service.Create(context.Background(), params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repository.createCalled {
		t.Fatal("expected repository to be called")
	}
}
