package trips

import (
	"context"

	"github.com/goretante/clearmytravel-api/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBRepository struct {
	queries *db.Queries
}

func NewDBRepository(queries *db.Queries) *DBRepository {
	return &DBRepository{
		queries: queries,
	}
}

func (r *DBRepository) Create(
	ctx context.Context,
	params db.CreateTripParams,
) (db.Trip, error) {
	return r.queries.CreateTrip(ctx, params)
}

func (r *DBRepository) GetByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.Trip, error) {
	return r.queries.GetTripByID(ctx, id)
}

func (r *DBRepository) ListByUserID(
	ctx context.Context,
	userID pgtype.UUID,
) ([]db.Trip, error) {
	return r.queries.ListTripsByUserID(ctx, userID)
}

func (r *DBRepository) ListByPassportID(
	ctx context.Context,
	passportID pgtype.UUID,
) ([]db.Trip, error) {
	return r.queries.ListTripsByPassportID(ctx, passportID)
}

func (r *DBRepository) Update(
	ctx context.Context,
	params db.UpdateTripParams,
) (db.Trip, error) {
	return r.queries.UpdateTrip(ctx, params)
}

func (r *DBRepository) Delete(
	ctx context.Context,
	id pgtype.UUID,
) error {
	return r.queries.DeleteTrip(ctx, id)
}
