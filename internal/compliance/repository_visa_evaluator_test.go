package compliance

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func (r fakeVisaRuleRepository) FindVisaRule(
	ctx context.Context,
	nationalityCode string,
	destinationCode string,
	departureDate time.Time,
) (VisaRule, error) {
	return r.rule, r.err
}

func TestRepositoryVisaEvaluatorRuleFound(t *testing.T) {
	evaluator := NewRepositoryVisaEvaluator(
		fakeVisaRuleRepository{
			rule: VisaRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				VisaRequired:    false,
			},
		},
	)

	result := evaluator.Evaluate(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		VisaStatus:      VisaUnknown,
	})

	if result.Status != RequirementOK {
		t.Fatalf("expected status %q, got %q", RequirementOK, result.Status)
	}

	if result.Message != "visa is not required" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}

func TestRepositoryVisaEvaluatorRuleFoundAndVisaMissing(t *testing.T) {
	evaluator := NewRepositoryVisaEvaluator(
		fakeVisaRuleRepository{
			rule: VisaRule{
				NationalityCode: "HRV",
				DestinationCode: "CHN",
				VisaRequired:    true,
			},
		},
	)

	result := evaluator.Evaluate(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "CHN",
		VisaStatus:      VisaMissing,
	})

	if result.Status != RequirementProblem {
		t.Fatalf("expected status %q, got %q", RequirementProblem, result.Status)
	}

	if result.Message != "visa is required but has not been obtained" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}

func TestRepositoryVisaEvaluatorNoRule(t *testing.T) {
	evaluator := NewRepositoryVisaEvaluator(
		fakeVisaRuleRepository{
			err: pgx.ErrNoRows,
		},
	)

	result := evaluator.Evaluate(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "XXX",
	})

	if result.Status != RequirementUnknown {
		t.Fatalf("expected status %q, got %q", RequirementUnknown, result.Status)
	}

	if result.Message != "visa requirement could not be determined" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}

func TestRepositoryVisaEvaluatorRepositoryError(t *testing.T) {
	evaluator := NewRepositoryVisaEvaluator(
		fakeVisaRuleRepository{
			err: errors.New("database error"),
		},
	)

	result := evaluator.Evaluate(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
	})

	if result.Status != RequirementUnknown {
		t.Fatalf("expected status %q, got %q", RequirementUnknown, result.Status)
	}

	if result.Message != "visa requirement could not be determined" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}
