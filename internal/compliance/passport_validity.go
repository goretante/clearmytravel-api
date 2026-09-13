package compliance

import "context"

type PassportValidityEvaluator struct {
	BufferDays int
}

func (e PassportValidityEvaluator) Evaluate(
	ctx context.Context,
	req CheckRequest,
) Requirement {
	requiredUntil := req.ReturnDate.AddDate(0, 0, e.BufferDays)

	requirement := Requirement{
		Type:   "passport_validity",
		Status: RequirementOK,
	}

	if !req.PassportExpiry.After(requiredUntil) {
		requirement.Status = RequirementProblem
		requirement.Message = "passport does not have sufficient validity after the trip"

		return requirement
	}

	requirement.Message = "passport has sufficient validity after the trip"
	return requirement
}
