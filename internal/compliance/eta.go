package compliance

import "time"

type ETAStatus string

const (
	ETANotRequired ETAStatus = "not_required"
	ETAObtained    ETAStatus = "obtained"
	ETAMissing     ETAStatus = "missing"
	ETAUnknown     ETAStatus = "unknown"
)

type ETARule struct {
	NationalityCode string
	DestinationCode string
	ETARequired     bool

	EffectiveFrom time.Time
	EffectiveTo   *time.Time

	Source RuleSource
}

func evaluateETARule(
	rule ETARule,
	req CheckRequest,
) Requirement {
	requirement := Requirement{
		Type:   "eta",
		Source: &rule.Source,
	}

	if !rule.ETARequired {
		requirement.Status = RequirementOK
		requirement.Message = "eta is not required"
		return requirement
	}

	switch req.ETAStatus {
	case ETAObtained:
		requirement.Status = RequirementOK
		requirement.Message = "eta has been obtained"

	case ETAMissing:
		requirement.Status = RequirementProblem
		requirement.Message = "eta is required but has not been obtained"

	default:
		requirement.Status = RequirementUnknown
		requirement.Message = "eta is required but eta status is unknown"
	}

	return requirement
}
