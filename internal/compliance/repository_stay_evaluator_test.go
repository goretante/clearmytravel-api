package compliance

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

type fakeStayRuleRepository struct {
	rule StayRule
	err  error
}

func (r fakeStayRuleRepository) FindStayRule(
	ctx context.Context,
	nationalityCode string,
	destinationCode string,
	departureDate time.Time,
) (StayRule, error) {
	return r.rule, r.err
}

func TestRepositoryStayDurationEvaluatorRuleFound(t *testing.T) {
	evaluator := NewRepositoryStayDurationEvaluator(
		fakeStayRuleRepository{
			rule: StayRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				MaxStayDays:     90,
			},
		},
	)

	result := evaluator.Evaluate(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2031, 7, 15, 0, 0, 0, 0, time.UTC),
	})

	if result.Status != RequirementOK {
		t.Fatalf("expected status %q, got %q", RequirementOK, result.Status)
	}

	if result.Message != "planned stay is within the allowed duration" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}

func TestRepositoryStayDurationEvaluatorRuleFoundAndStayExceeded(t *testing.T) {
	evaluator := NewRepositoryStayDurationEvaluator(
		fakeStayRuleRepository{
			rule: StayRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				MaxStayDays:     30,
			},
		},
	)

	result := evaluator.Evaluate(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2031, 7, 15, 0, 0, 0, 0, time.UTC),
	})

	if result.Status != RequirementProblem {
		t.Fatalf("expected status %q, got %q", RequirementProblem, result.Status)
	}

	if result.Message != "planned stay exceeds the maximum allowed duration" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}

func TestRepositoryStayDurationEvaluatorNoRule(t *testing.T) {
	evaluator := NewRepositoryStayDurationEvaluator(
		fakeStayRuleRepository{
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

	if result.Message != "stay duration requirement could not be determined" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}

func TestRepositoryStayDurationEvaluatorRepositoryError(t *testing.T) {
	evaluator := NewRepositoryStayDurationEvaluator(
		fakeStayRuleRepository{
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

	if result.Message != "stay duration requirement could not be determined" {
		t.Fatalf("unexpected message: %q", result.Message)
	}
}
