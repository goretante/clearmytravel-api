package compliance

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/goretante/clearmytravel-api/internal/db"
)

type DBETARuleRepository struct {
	queries *db.Queries
}

func NewDBETARuleRepository(queries *db.Queries) *DBETARuleRepository {
	return &DBETARuleRepository{
		queries: queries,
	}
}

func (r *DBETARuleRepository) FindETARule(
	ctx context.Context,
	nationalityCode string,
	destinationCode string,
	departureDate time.Time,
) (ETARule, error) {
	row, err := r.queries.FindETARule(
		ctx,
		db.FindETARuleParams{
			NationalityCode: nationalityCode,
			DestinationCode: destinationCode,
			EffectiveFrom: pgtype.Date{
				Time:  departureDate,
				Valid: true,
			},
		},
	)
	if err != nil {
		return ETARule{}, err
	}

	return ETARule{
		NationalityCode: row.NationalityCode,
		DestinationCode: row.DestinationCode,
		ETARequired:     row.EtaRequired,
		EffectiveFrom:   row.EffectiveFrom.Time,
		EffectiveTo:     datePtr(row.EffectiveTo),
		Source: RuleSource{
			Name: row.SourceName,
			URL:  row.SourceUrl,
		},
	}, nil
}
