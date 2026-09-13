package compliance

import (
	"context"
	"time"

	"github.com/goretante/clearmytravel-api/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBVisaRuleRepository struct {
	queries *db.Queries
}

func NewDBVisaRuleRepository(queries *db.Queries) *DBVisaRuleRepository {
	return &DBVisaRuleRepository{
		queries: queries,
	}
}

func (r *DBVisaRuleRepository) FindVisaRule(
	ctx context.Context,
	nationalityCode string,
	destinationCode string,
	departureDate time.Time,
) (VisaRule, error) {
	row, err := r.queries.FindVisaRule(
		ctx,
		db.FindVisaRuleParams{
			NationalityCode: nationalityCode,
			DestinationCode: destinationCode,
			EffectiveFrom: pgtype.Date{
				Time:  departureDate,
				Valid: true,
			},
		},
	)
	if err != nil {
		return VisaRule{}, err
	}

	return VisaRule{
		NationalityCode: row.NationalityCode,
		DestinationCode: row.DestinationCode,
		VisaRequired:    row.VisaRequired,
		EffectiveFrom:   row.EffectiveFrom.Time,
		EffectiveTo:     datePtr(row.EffectiveTo),
		Source: RuleSource{
			Name: row.SourceName,
			URL:  row.SourceUrl,
		},
	}, nil
}

func datePtr(value pgtype.Date) *time.Time {
	if !value.Valid {
		return nil
	}

	t := value.Time
	return &t
}
