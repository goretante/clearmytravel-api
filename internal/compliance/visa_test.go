package compliance

import "testing"

func TestEvaluateVisaRuleVisaNotRequired(t *testing.T) {
	rule := VisaRule{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		VisaRequired:    false,
	}

	result := evaluateVisaRule(rule, CheckRequest{
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

	if result.Message != "visa is not required" {
		t.Fatalf(
			"expected message %q, got %q",
			"visa is not required",
			result.Message,
		)
	}
}

func TestEvaluateVisaRuleVisaRequiredButStatusUnknown(t *testing.T) {
	rule := VisaRule{
		NationalityCode: "HRV",
		DestinationCode: "CHN",
		VisaRequired:    true,
	}

	result := evaluateVisaRule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "CHN",
	})

	if result.Status != RequirementUnknown {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementUnknown,
			result.Status,
		)
	}

	if result.Message != "visa is required but visa status is unknown" {
		t.Fatalf(
			"expected message %q, got %q",
			"visa is required but visa status is unknown",
			result.Message,
		)
	}
}

func TestEvaluateVisaRuleVisaObtained(t *testing.T) {
	rule := VisaRule{
		NationalityCode: "HRV",
		DestinationCode: "CHN",
		VisaRequired:    true,
	}

	result := evaluateVisaRule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "CHN",
		VisaStatus:      VisaObtained,
	})

	if result.Status != RequirementOK {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementOK,
			result.Status,
		)
	}

	if result.Message != "visa has been obtained" {
		t.Fatalf(
			"expected message %q, got %q",
			"visa has been obtained",
			result.Message,
		)
	}
}

func TestEvaluateVisaRuleVisaMissing(t *testing.T) {
	rule := VisaRule{
		NationalityCode: "HRV",
		DestinationCode: "CHN",
		VisaRequired:    true,
	}

	result := evaluateVisaRule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "CHN",
		VisaStatus:      VisaMissing,
	})

	if result.Status != RequirementProblem {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementProblem,
			result.Status,
		)
	}

	if result.Message != "visa is required but has not been obtained" {
		t.Fatalf(
			"expected message %q, got %q",
			"visa is required but has not been obtained",
			result.Message,
		)
	}
}

func TestEvaluateVisaRuleVisaUnknown(t *testing.T) {
	rule := VisaRule{
		NationalityCode: "HRV",
		DestinationCode: "CHN",
		VisaRequired:    true,
	}

	result := evaluateVisaRule(rule, CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "CHN",
		VisaStatus:      VisaUnknown,
	})

	if result.Status != RequirementUnknown {
		t.Fatalf(
			"expected status %q, got %q",
			RequirementUnknown,
			result.Status,
		)
	}

	if result.Message != "visa is required but visa status is unknown" {
		t.Fatalf(
			"expected message %q, got %q",
			"visa is required but visa status is unknown",
			result.Message,
		)
	}
}
