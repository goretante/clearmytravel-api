package compliance

import (
	"context"
)

type Service struct {
	evaluators []Evaluator
}

func NewService(evaluators ...Evaluator) *Service {
	return &Service{
		evaluators: evaluators,
	}
}

func (s *Service) Check(ctx context.Context, req CheckRequest) CheckResult {
	result := CheckResult{
		Status:       StatusCleared,
		Requirements: make([]Requirement, 0, len(s.evaluators)),
	}

	for _, evaluator := range s.evaluators {
		requirement := evaluator.Evaluate(ctx, req)

		result.Requirements = append(
			result.Requirements,
			requirement,
		)

		if requirement.Status == RequirementProblem ||
			requirement.Status == RequirementUnknown {
			result.Status = StatusActionRequired
		}
	}

	return result
}
