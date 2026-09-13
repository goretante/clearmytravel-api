package compliance

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type RepositoryStayDurationEvaluator struct {
	repository StayRuleRepository
}

func NewRepositoryStayDurationEvaluator(
	repository StayRuleRepository,
) RepositoryStayDurationEvaluator {
	return RepositoryStayDurationEvaluator{
		repository: repository,
	}
}

func (e RepositoryStayDurationEvaluator) Evaluate(
	ctx context.Context,
	req CheckRequest,
) Requirement {
	rule, err := e.repository.FindStayRule(
		ctx,
		req.NationalityCode,
		req.DestinationCode,
		req.DepartureDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Requirement{
				Type:    "stay_duration",
				Status:  RequirementUnknown,
				Message: "stay duration requirement could not be determined",
			}
		}

		return Requirement{
			Type:    "stay_duration",
			Status:  RequirementUnknown,
			Message: "stay duration requirement could not be determined",
		}
	}

	return evaluateStayRule(rule, req)
}
