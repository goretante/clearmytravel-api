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

func TestCalculateStayDays(t *testing.T) {
	tests := []struct {
		name       string
		departure  time.Time
		returnDate time.Time
		expected   int
	}{
		{
			name: "same day",
			departure: time.Date(
				2031, 6, 1, 10, 0, 0, 0, time.UTC,
			),
			returnDate: time.Date(
				2031, 6, 1, 18, 0, 0, 0, time.UTC,
			),
			expected: 0,
		},
		{
			name: "one calendar day",
			departure: time.Date(
				2031, 6, 1, 23, 30, 0, 0, time.UTC,
			),
			returnDate: time.Date(
				2031, 6, 2, 8, 0, 0, 0, time.UTC,
			),
			expected: 1,
		},
		{
			name: "nineteen calendar days",
			departure: time.Date(
				2031, 6, 1, 23, 30, 0, 0, time.UTC,
			),
			returnDate: time.Date(
				2031, 6, 20, 8, 0, 0, 0, time.UTC,
			),
			expected: 19,
		},
		{
			name: "ninety calendar days",
			departure: time.Date(
				2031, 1, 1, 23, 30, 0, 0, time.UTC,
			),
			returnDate: time.Date(
				2031, 4, 1, 8, 0, 0, 0, time.UTC,
			),
			expected: 90,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateStayDays(
				tt.departure,
				tt.returnDate,
			)

			if got != tt.expected {
				t.Fatalf(
					"expected %d stay days, got %d",
					tt.expected,
					got,
				)
			}
		})
	}
}

func TestEvaluateStayRuleExactLimitWithDifferentTimes(t *testing.T) {
	rule := StayRule{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		MaxStayDays:     90,
	}

	result := evaluateStayRule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate: time.Date(
			2031, 1, 1, 23, 30, 0, 0, time.UTC,
		),
		ReturnDate: time.Date(
			2031, 4, 1, 8, 0, 0, 0, time.UTC,
		),
	})

	if result.Status != RequirementOK {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementOK,
			result.Status,
		)
	}
}

func TestEvaluateStayRuleExceedsLimitWithDifferentTimes(t *testing.T) {
	rule := StayRule{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		MaxStayDays:     90,
	}

	result := evaluateStayRule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate: time.Date(
			2031, 1, 1, 23, 30, 0, 0, time.UTC,
		),
		ReturnDate: time.Date(
			2031, 4, 2, 8, 0, 0, 0, time.UTC,
		),
	})

	if result.Status != RequirementProblem {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementProblem,
			result.Status,
		)
	}
}
