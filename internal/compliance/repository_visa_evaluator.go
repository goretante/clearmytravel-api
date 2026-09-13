package compliance

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type RepositoryVisaEvaluator struct {
	repository VisaRuleRepository
}

func NewRepositoryVisaEvaluator(
	repository VisaRuleRepository,
) RepositoryVisaEvaluator {
	return RepositoryVisaEvaluator{
		repository: repository,
	}
}

func (e RepositoryVisaEvaluator) Evaluate(
	ctx context.Context,
	req CheckRequest,
) Requirement {
	rule, err := e.repository.FindVisaRule(
		ctx,
		req.NationalityCode,
		req.DestinationCode,
		req.DepartureDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Requirement{
				Type:    "visa",
				Status:  RequirementUnknown,
				Message: "visa requirement could not be determined",
			}
		}

		return Requirement{
			Type:    "visa",
			Status:  RequirementUnknown,
			Message: "visa requirement could not be determined",
		}
	}

	return evaluateVisaRule(rule, req)
}
