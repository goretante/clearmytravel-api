package compliance

import (
	"context"
	"testing"
	"time"
)

type fakeEvaluator struct {
	requirement Requirement
}

func (e fakeEvaluator) Evaluate(
	ctx context.Context,
	req CheckRequest,
) Requirement {
	return e.requirement
}

type fakeVisaRuleRepository struct {
	rule VisaRule
	err  error
}

func TestServiceCheck(t *testing.T) {
	returnDate := time.Date(2031, 6, 10, 0, 0, 0, 0, time.UTC)

	service := NewService(
		PassportValidityEvaluator{
			BufferDays: 0,
		},
	)

	result := service.Check(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      returnDate,
		PassportExpiry:  time.Date(2031, 6, 15, 0, 0, 0, 0, time.UTC),
	})

	if result.Status != StatusCleared {
		t.Fatalf("expected status %q, got %q", StatusCleared, result.Status)
	}

	if len(result.Requirements) != 1 {
		t.Fatalf("expected 1 requirement, got %d", len(result.Requirements))
	}

	requirement := result.Requirements[0]

	if requirement.Type != "passport_validity" {
		t.Fatalf(
			"expected requirement type %q, got %q",
			"passport_validity",
			requirement.Type,
		)
	}

	if requirement.Status != RequirementOK {
		t.Fatalf(
			"expected requirement status %q, got %q",
			RequirementOK,
			requirement.Status,
		)
	}
}

func TestServiceCheckPassportExpiredBeforeReturn(t *testing.T) {
	service := NewService(
		PassportValidityEvaluator{
			BufferDays: 0,
		},
	)

	result := service.Check(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2031, 6, 10, 0, 0, 0, 0, time.UTC),
		PassportExpiry:  time.Date(2031, 6, 5, 0, 0, 0, 0, time.UTC),
	})

	if result.Status != StatusActionRequired {
		t.Fatalf(
			"expected status %q, got %q",
			StatusActionRequired,
			result.Status,
		)
	}

	if len(result.Requirements) != 1 {
		t.Fatalf("expected 1 requirement, got %d", len(result.Requirements))
	}

	if result.Requirements[0].Status != RequirementProblem {
		t.Fatalf(
			"expected requirement status %q, got %q",
			RequirementProblem,
			result.Requirements[0].Status,
		)
	}
}

func TestServiceCheckMultipleEvaluators(t *testing.T) {
	service := NewService(
		fakeEvaluator{
			requirement: Requirement{
				Type:    "passport_validity",
				Status:  RequirementOK,
				Message: "passport is valid",
			},
		},
		fakeEvaluator{
			requirement: Requirement{
				Type:    "visa",
				Status:  RequirementProblem,
				Message: "visa is required",
			},
		},
	)

	result := service.Check(
		context.Background(),
		CheckRequest{},
	)

	if result.Status != StatusActionRequired {
		t.Fatalf(
			"expected status %q, got %q",
			StatusActionRequired,
			result.Status,
		)
	}

	if len(result.Requirements) != 2 {
		t.Fatalf(
			"expected 2 requirements, got %d",
			len(result.Requirements),
		)
	}

	if result.Requirements[0].Status != RequirementOK {
		t.Fatalf("expected first requirement to be ok")
	}

	if result.Requirements[1].Status != RequirementProblem {
		t.Fatalf("expected second requirement to be problem")
	}
}

func TestServiceCheckUnknownRequirement(t *testing.T) {
	service := NewService(
		fakeEvaluator{
			requirement: Requirement{
				Type:    "visa",
				Status:  RequirementUnknown,
				Message: "visa requirement could not be determined",
			},
		},
	)

	result := service.Check(
		context.Background(),
		CheckRequest{},
	)

	if result.Status != StatusActionRequired {
		t.Fatalf(
			"expected status %q, got %q",
			StatusActionRequired,
			result.Status,
		)
	}

	if len(result.Requirements) != 1 {
		t.Fatalf(
			"expected 1 requirement, got %d",
			len(result.Requirements),
		)
	}

	if result.Requirements[0].Status != RequirementUnknown {
		t.Fatalf(
			"expected requirement status %q, got %q",
			RequirementUnknown,
			result.Requirements[0].Status,
		)
	}
}

func TestServiceCheckPassportAndVisa(t *testing.T) {
	visaEvaluator := NewRepositoryVisaEvaluator(
		fakeVisaRuleRepository{
			rule: VisaRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				VisaRequired:    false,
			},
		},
	)

	service := NewService(
		PassportValidityEvaluator{
			BufferDays: 0,
		},
		visaEvaluator,
	)

	result := service.Check(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2031, 6, 10, 0, 0, 0, 0, time.UTC),
		PassportExpiry:  time.Date(2031, 6, 15, 0, 0, 0, 0, time.UTC),
	})

	if result.Status != StatusCleared {
		t.Fatalf(
			"expected status %q, got %q",
			StatusCleared,
			result.Status,
		)
	}

	if len(result.Requirements) != 2 {
		t.Fatalf(
			"expected 2 requirements, got %d",
			len(result.Requirements),
		)
	}

	if result.Requirements[0].Status != RequirementOK {
		t.Fatalf("expected passport requirement to be ok")
	}

	if result.Requirements[1].Status != RequirementOK {
		t.Fatalf("expected visa requirement to be ok")
	}
}

func TestServiceCheckVisaRequired(t *testing.T) {
	visaEvaluator := NewRepositoryVisaEvaluator(
		fakeVisaRuleRepository{
			rule: VisaRule{
				NationalityCode: "HRV",
				DestinationCode: "CHN",
				VisaRequired:    true,
			},
		},
	)

	service := NewService(
		PassportValidityEvaluator{
			BufferDays: 0,
		},
		visaEvaluator,
	)

	result := service.Check(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "CHN",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2031, 6, 10, 0, 0, 0, 0, time.UTC),
		PassportExpiry:  time.Date(2031, 6, 15, 0, 0, 0, 0, time.UTC),
		VisaStatus:      VisaMissing,
	})

	if result.Status != StatusActionRequired {
		t.Fatalf(
			"expected status %q, got %q",
			StatusActionRequired,
			result.Status,
		)
	}

	if len(result.Requirements) != 2 {
		t.Fatalf(
			"expected 2 requirements, got %d",
			len(result.Requirements),
		)
	}

	if result.Requirements[0].Status != RequirementOK {
		t.Fatalf("expected passport requirement to be ok")
	}

	if result.Requirements[1].Status != RequirementProblem {
		t.Fatalf("expected visa requirement to be problem")
	}
}

func TestServiceCheckPassportVisaAndETA(t *testing.T) {
	visaEvaluator := NewRepositoryVisaEvaluator(
		fakeVisaRuleRepository{
			rule: VisaRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				VisaRequired:    false,
			},
		},
	)

	etaEvaluator := NewRepositoryETAEvaluator(
		fakeETARuleRepository{
			rule: ETARule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				ETARequired:     true,
			},
		},
	)

	service := NewService(
		PassportValidityEvaluator{
			BufferDays: 0,
		},
		visaEvaluator,
		etaEvaluator,
	)

	result := service.Check(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   time.Date(2031, 6, 1, 0, 0, 0, 0, time.UTC),
		ReturnDate:      time.Date(2031, 6, 10, 0, 0, 0, 0, time.UTC),
		PassportExpiry:  time.Date(2031, 6, 15, 0, 0, 0, 0, time.UTC),
		ETAStatus:       ETAMissing,
	})

	if result.Status != StatusActionRequired {
		t.Fatalf(
			"expected status %q, got %q",
			StatusActionRequired,
			result.Status,
		)
	}

	if len(result.Requirements) != 3 {
		t.Fatalf(
			"expected 3 requirements, got %d",
			len(result.Requirements),
		)
	}

	if result.Requirements[0].Type != "passport_validity" {
		t.Fatalf(
			"expected first requirement type %q, got %q",
			"passport_validity",
			result.Requirements[0].Type,
		)
	}

	if result.Requirements[0].Status != RequirementOK {
		t.Fatalf("expected passport requirement to be ok")
	}

	if result.Requirements[1].Type != "visa" {
		t.Fatalf(
			"expected second requirement type %q, got %q",
			"visa",
			result.Requirements[1].Type,
		)
	}

	if result.Requirements[1].Status != RequirementOK {
		t.Fatalf("expected visa requirement to be ok")
	}

	if result.Requirements[2].Type != "eta" {
		t.Fatalf(
			"expected third requirement type %q, got %q",
			"eta",
			result.Requirements[2].Type,
		)
	}

	if result.Requirements[2].Status != RequirementProblem {
		t.Fatalf("expected eta requirement to be problem")
	}
}

func TestServiceCheckAllRequirementsOK(t *testing.T) {
	departure := time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC)
	returnDate := time.Date(2027, 6, 15, 0, 0, 0, 0, time.UTC)
	passportExpiry := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	visaEvaluator := NewRepositoryVisaEvaluator(
		fakeVisaRuleRepository{
			rule: VisaRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				VisaRequired:    false,
			},
		},
	)

	etaEvaluator := NewRepositoryETAEvaluator(
		fakeETARuleRepository{
			rule: ETARule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				ETARequired:     true,
			},
		},
	)

	stayEvaluator := NewRepositoryStayDurationEvaluator(
		fakeStayRuleRepository{
			rule: StayRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				MaxStayDays:     90,
			},
		},
	)

	service := NewService(
		PassportValidityEvaluator{BufferDays: 90},
		visaEvaluator,
		etaEvaluator,
		stayEvaluator,
	)

	result := service.Check(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   departure,
		ReturnDate:      returnDate,
		PassportExpiry:  passportExpiry,
		VisaStatus:      VisaNotRequired,
		ETAStatus:       ETAObtained,
	})

	if result.Status != StatusCleared {
		t.Fatalf(
			"expected status %q, got %q",
			StatusCleared,
			result.Status,
		)
	}

	if len(result.Requirements) != 4 {
		t.Fatalf(
			"expected 4 requirements, got %d",
			len(result.Requirements),
		)
	}

	for i, requirement := range result.Requirements {
		if requirement.Status != RequirementOK {
			t.Errorf(
				"requirement %d (%s): expected status %q, got %q",
				i,
				requirement.Type,
				RequirementOK,
				requirement.Status,
			)
		}
	}
}

func TestServiceCheckAllRequirementsWithMissingETA(t *testing.T) {
	departure := time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC)
	returnDate := time.Date(2027, 6, 15, 0, 0, 0, 0, time.UTC)
	passportExpiry := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	visaEvaluator := NewRepositoryVisaEvaluator(
		fakeVisaRuleRepository{
			rule: VisaRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				VisaRequired:    false,
			},
		},
	)

	etaEvaluator := NewRepositoryETAEvaluator(
		fakeETARuleRepository{
			rule: ETARule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				ETARequired:     true,
			},
		},
	)

	stayEvaluator := NewRepositoryStayDurationEvaluator(
		fakeStayRuleRepository{
			rule: StayRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				MaxStayDays:     90,
			},
		},
	)

	service := NewService(
		PassportValidityEvaluator{BufferDays: 90},
		visaEvaluator,
		etaEvaluator,
		stayEvaluator,
	)

	result := service.Check(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   departure,
		ReturnDate:      returnDate,
		PassportExpiry:  passportExpiry,
		VisaStatus:      VisaNotRequired,
		ETAStatus:       ETAMissing,
	})

	if result.Status != StatusActionRequired {
		t.Fatalf(
			"expected status %q, got %q",
			StatusActionRequired,
			result.Status,
		)
	}

	if len(result.Requirements) != 4 {
		t.Fatalf(
			"expected 4 requirements, got %d",
			len(result.Requirements),
		)
	}

	expectedStatuses := []RequirementStatus{
		RequirementOK,
		RequirementOK,
		RequirementProblem,
		RequirementOK,
	}

	for i, expected := range expectedStatuses {
		if result.Requirements[i].Status != expected {
			t.Errorf(
				"requirement %d (%s): expected status %q, got %q",
				i,
				result.Requirements[i].Type,
				expected,
				result.Requirements[i].Status,
			)
		}
	}
}

func TestServiceCheckStayDurationExceeded(t *testing.T) {
	departure := time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC)
	returnDate := time.Date(2027, 10, 1, 0, 0, 0, 0, time.UTC)
	passportExpiry := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	visaEvaluator := NewRepositoryVisaEvaluator(
		fakeVisaRuleRepository{
			rule: VisaRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				VisaRequired:    false,
			},
		},
	)

	etaEvaluator := NewRepositoryETAEvaluator(
		fakeETARuleRepository{
			rule: ETARule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				ETARequired:     true,
			},
		},
	)

	stayEvaluator := NewRepositoryStayDurationEvaluator(
		fakeStayRuleRepository{
			rule: StayRule{
				NationalityCode: "HRV",
				DestinationCode: "USA",
				MaxStayDays:     90,
			},
		},
	)

	service := NewService(
		PassportValidityEvaluator{BufferDays: 90},
		visaEvaluator,
		etaEvaluator,
		stayEvaluator,
	)

	result := service.Check(context.Background(), CheckRequest{
		NationalityCode: "HRV",
		DestinationCode: "USA",
		DepartureDate:   departure,
		ReturnDate:      returnDate,
		PassportExpiry:  passportExpiry,
		VisaStatus:      VisaNotRequired,
		ETAStatus:       ETAObtained,
	})

	if result.Status != StatusActionRequired {
		t.Fatalf(
			"expected status %q, got %q",
			StatusActionRequired,
			result.Status,
		)
	}

	if len(result.Requirements) != 4 {
		t.Fatalf(
			"expected 4 requirements, got %d",
			len(result.Requirements),
		)
	}

	expectedStatuses := []RequirementStatus{
		RequirementOK,
		RequirementOK,
		RequirementOK,
		RequirementProblem,
	}

	for i, expected := range expectedStatuses {
		if result.Requirements[i].Status != expected {
			t.Errorf(
				"requirement %d (%s): expected status %q, got %q",
				i,
				result.Requirements[i].Type,
				expected,
				result.Requirements[i].Status,
			)
		}
	}
}
