package compliance

import "time"

type StayRule struct {
	NationalityCode string
	DestinationCode string
	MaxStayDays     int
	EffectiveFrom   time.Time
	EffectiveTo     *time.Time
	Source          RuleSource
}

func calculateStayDays(
	departure time.Time,
	returnDate time.Time,
) int {
	departureDate := dateOnly(departure)
	returnDateOnly := dateOnly(returnDate)

	return int(
		returnDateOnly.Sub(departureDate).Hours() / 24,
	)
}

func dateOnly(t time.Time) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func evaluateStayRule(
	rule StayRule,
	req CheckRequest,
) Requirement {
	requirement := Requirement{
		Type:   "stay_duration",
		Source: &rule.Source,
	}

	stayDays := calculateStayDays(
		req.DepartureDate,
		req.ReturnDate,
	)

	if stayDays > rule.MaxStayDays {
		requirement.Status = RequirementProblem
		requirement.Message = "planned stay exceeds the maximum allowed duration"
		return requirement
	}

	requirement.Status = RequirementOK
	requirement.Message = "planned stay is within the allowed duration"

	return requirement
}
