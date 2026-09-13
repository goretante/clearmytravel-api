package compliance

import "time"

type VisaStatus string

const (
	VisaNotRequired VisaStatus = "not_required"
	VisaRequired    VisaStatus = "required"
	VisaObtained    VisaStatus = "obtained"
	VisaMissing     VisaStatus = "missing"
	VisaUnknown     VisaStatus = "unknown"
)

type RuleSource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VisaRule struct {
	NationalityCode string
	DestinationCode string
	VisaRequired    bool

	EffectiveFrom time.Time
	EffectiveTo   *time.Time

	Source RuleSource
}

func evaluateVisaRule(
	rule VisaRule,
	req CheckRequest,
) Requirement {
	requirement := Requirement{
		Type:   "visa",
		Source: &rule.Source,
	}

	if !rule.VisaRequired {
		requirement.Status = RequirementOK
		requirement.Message = "visa is not required"
		return requirement
	}

	switch req.VisaStatus {
	case VisaObtained:
		requirement.Status = RequirementOK
		requirement.Message = "visa has been obtained"

	case VisaMissing:
		requirement.Status = RequirementProblem
		requirement.Message = "visa is required but has not been obtained"

	case VisaUnknown:
		requirement.Status = RequirementUnknown
		requirement.Message = "visa is required but visa status is unknown"

	default:
		requirement.Status = RequirementUnknown
		requirement.Message = "visa is required but visa status is unknown"
	}

	return requirement
}
