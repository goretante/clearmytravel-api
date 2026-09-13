package compliance

import "testing"

func TestEvaluateETARuleNotRequired(t *testing.T) {
	rule := ETARule{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		ETARequired:     false,
	}

	result := evaluateETARule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
	})

	if result.Status != RequirementOK {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementOK,
			result.Status,
		)
	}

	if result.Message != "eta is not required" {
		t.Fatalf(
			"expected message %q, got %q",
			"eta is not required",
			result.Message,
		)
	}
}

func TestEvaluateETARuleObtained(t *testing.T) {
	rule := ETARule{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		ETARequired:     true,
	}

	result := evaluateETARule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		ETAStatus:       ETAObtained,
	})

	if result.Status != RequirementOK {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementOK,
			result.Status,
		)
	}

	if result.Message != "eta has been obtained" {
		t.Fatalf(
			"expected message %q, got %q",
			"eta has been obtained",
			result.Message,
		)
	}
}

func TestEvaluateETARuleMissing(t *testing.T) {
	rule := ETARule{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		ETARequired:     true,
	}

	result := evaluateETARule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		ETAStatus:       ETAMissing,
	})

	if result.Status != RequirementProblem {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementProblem,
			result.Status,
		)
	}

	if result.Message != "eta is required but has not been obtained" {
		t.Fatalf(
			"expected message %q, got %q",
			"eta is required but has not been obtained",
			result.Message,
		)
	}
}

func TestEvaluateETARuleUnknown(t *testing.T) {
	rule := ETARule{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		ETARequired:     true,
	}

	result := evaluateETARule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		ETAStatus:       ETAUnknown,
	})

	if result.Status != RequirementUnknown {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementUnknown,
			result.Status,
		)
	}

	if result.Message != "eta is required but eta status is unknown" {
		t.Fatalf(
			"expected message %q, got %q",
			"eta is required but eta status is unknown",
			result.Message,
		)
	}
}
