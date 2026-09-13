package compliance

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type RepositoryETAEvaluator struct {
	repository ETARuleRepository
}

func NewRepositoryETAEvaluator(
	repository ETARuleRepository,
) RepositoryETAEvaluator {
	return RepositoryETAEvaluator{
		repository: repository,
	}
}

func (e RepositoryETAEvaluator) Evaluate(
	ctx context.Context,
	req CheckRequest,
) Requirement {
	rule, err := e.repository.FindETARule(
		ctx,
		req.NationalityCode,
		req.DestinationCode,
		req.DepartureDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Requirement{
				Type:    "eta",
				Status:  RequirementUnknown,
				Message: "eta requirement could not be determined",
			}
		}

		return Requirement{
			Type:    "eta",
			Status:  RequirementUnknown,
			Message: "eta requirement could not be determined",
		}
	}

	return evaluateETARule(rule, req)
}
