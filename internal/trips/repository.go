package trips

import (
	"context"

	"github.com/goretante/clearmytravel-api/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	Create(
		ctx context.Context,
		params db.CreateTripParams,
	) (db.Trip, error)

	GetByID(
		ctx context.Context,
		id pgtype.UUID,
	) (db.Trip, error)

	ListByUserID(
		ctx context.Context,
		userID pgtype.UUID,
	) ([]db.Trip, error)

	ListByPassportID(
		ctx context.Context,
		passportID pgtype.UUID,
	) ([]db.Trip, error)

	Update(
		ctx context.Context,
		params db.UpdateTripParams,
	) (db.Trip, error)

	Delete(
		ctx context.Context,
		id pgtype.UUID,
	) error
}
