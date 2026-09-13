package trips

import (
	"context"
	"errors"

	"github.com/goretante/clearmytravel-api/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidDateRange = errors.New(
		"return_date must be on or after departure_date",
	)
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	params db.CreateTripParams,
) (db.Trip, error) {
	if params.ReturnDate.Time.Before(params.DepartureDate.Time) {
		return db.Trip{}, ErrInvalidDateRange
	}

	return s.repository.Create(ctx, params)
}

func (s *Service) GetByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.Trip, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) ListByUserID(
	ctx context.Context,
	userID pgtype.UUID,
) ([]db.Trip, error) {
	return s.repository.ListByUserID(ctx, userID)
}

func (s *Service) ListByPassportID(
	ctx context.Context,
	passportID pgtype.UUID,
) ([]db.Trip, error) {
	return s.repository.ListByPassportID(ctx, passportID)
}

func (s *Service) Update(
	ctx context.Context,
	params db.UpdateTripParams,
) (db.Trip, error) {
	if params.ReturnDate.Time.Before(params.DepartureDate.Time) {
		return db.Trip{}, ErrInvalidDateRange
	}

	return s.repository.Update(ctx, params)
}

func (s *Service) Delete(
	ctx context.Context,
	id pgtype.UUID,
) error {
	return s.repository.Delete(ctx, id)
}
