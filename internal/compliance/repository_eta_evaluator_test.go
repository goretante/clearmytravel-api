package compliance

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

type fakeETARuleRepository struct {
	rule ETARule
	err  error
}

func (r fakeETARuleRepository) FindETARule(
	ctx context.Context,
	nationalityCode string,
	destinationCode string,
	departureDate time.Time,
) (ETARule, error) {
	return r.rule, r.err
}

func TestRepositoryETAEvaluatorRuleFound(t *testing.T) {
	evaluator := NewRepositoryETAEvaluator(
		fakeETARuleRepository{
			rule: ETARule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				ETARequired:     false,
			},
		},
	)

	result := evaluator.Evaluate(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		ETAStatus:       ETAUnknown,
	})

	if result.Status != RequirementOK {
		t.Fatalf("expected status %q, got %q", RequirementOK, result.Status)
	}

	if result.Message != "eta is not required" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}

func TestRepositoryETAEvaluatorRuleFoundAndETAMissing(t *testing.T) {
	evaluator := NewRepositoryETAEvaluator(
		fakeETARuleRepository{
			rule: ETARule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				ETARequired:     true,
			},
		},
	)

	result := evaluator.Evaluate(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		ETAStatus:       ETAMissing,
	})

	if result.Status != RequirementProblem {
		t.Fatalf("expected status %q, got %q", RequirementProblem, result.Status)
	}

	if result.Message != "eta is required but has not been obtained" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}

func TestRepositoryETAEvaluatorNoRule(t *testing.T) {
	evaluator := NewRepositoryETAEvaluator(
		fakeETARuleRepository{
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

	if result.Message != "eta requirement could not be determined" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}

func TestRepositoryETAEvaluatorRepositoryError(t *testing.T) {
	evaluator := NewRepositoryETAEvaluator(
		fakeETARuleRepository{
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

	if result.Message != "eta requirement could not be determined" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}
