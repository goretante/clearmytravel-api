package compliance

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/goretante/clearmytravel-api/internal/db"
)

type DBStayRuleRepository struct {
	queries *db.Queries
}

func NewDBStayRuleRepository(queries *db.Queries) *DBStayRuleRepository {
	return &DBStayRuleRepository{
		queries: queries,
	}
}

func (r *DBStayRuleRepository) FindStayRule(
	ctx context.Context,
	nationalityCode string,
	destinationCode string,
	departureDate time.Time,
) (StayRule, error) {
	row, err := r.queries.FindStayRule(
		ctx,
		db.FindStayRuleParams{
			NationalityCode: nationalityCode,
			DestinationCode: destinationCode,
			EffectiveFrom: pgtype.Date{
				Time:  departureDate,
				Valid: true,
			},
		},
	)
	if err != nil {
		return StayRule{}, err
	}

	return StayRule{
		NationalityCode: row.NationalityCode,
		DestinationCode: row.DestinationCode,
		MaxStayDays:     int(row.MaxStayDays),
		EffectiveFrom:   row.EffectiveFrom.Time,
		EffectiveTo:     datePtr(row.EffectiveTo),
		Source: RuleSource{
			Name: row.SourceName,
			URL:  row.SourceUrl,
		},
	}, nil
}
