package compliance

import (
	"testing"
	"time"
)

func TestEvaluateStayRuleWithinLimit(t *testing.T) {
	rule := StayRule{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		MaxStayDays:     90,
	}

	result := evaluateStayRule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2031, 7, 15, 0, 0, 0, 0, time.UTC),
	})

	if result.Status != RequirementOK {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementOK,
			result.Status,
		)
	}

	if result.Message != "planned stay is within the allowed duration" {
		t.Fatalf(
			"expected message %q, got %q",
			"planned stay is within the allowed duration",
			result.Message,
		)
	}
}

func TestEvaluateStayRuleExceedsLimit(t *testing.T) {
	rule := StayRule{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		MaxStayDays:     90,
	}

	result := evaluateStayRule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2031, 10, 1, 0, 0, 0, 0, time.UTC),
	})

	if result.Status != RequirementProblem {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementProblem,
			result.Status,
		)
	}

	if result.Message != "planned stay exceeds the maximum allowed duration" {
		t.Fatalf(
			"expected message %q, got %q",
			"planned stay exceeds the maximum allowed duration",
			result.Message,
		)
	}
}

func TestEvaluateStayRuleExactLimit(t *testing.T) {
	rule := StayRule{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		MaxStayDays:     90,
	}

	result := evaluateStayRule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2031, 8, 30, 0, 0, 0, 0, time.UTC),
	})

	if result.Status != RequirementOK {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementOK,
			result.Status,
		)
	}
}
