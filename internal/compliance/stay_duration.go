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

func evaluateStayRule(
	rule StayRule,
	req CheckRequest,
) Requirement {
	requirement := Requirement{
		Type:   "stay_duration",
		Source: &rule.Source,
	}

	stayDays := int(
		req.ReturnDate.Sub(req.DepartureDate).Hours() / 24,
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
