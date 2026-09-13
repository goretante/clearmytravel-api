package compliance

import "time"

type CheckRequest struct {
	NationalityCode string
	DestinationCode string
	DepartureDate   time.Time
	ReturnDate      time.Time
	PassportExpiry  time.Time
	VisaStatus      VisaStatus
	ETAStatus       ETAStatus
}

type Status string

const (
	StatusCleared        Status = "cleared"
	StatusActionRequired Status = "action_required"
)

type RequirementStatus string

const (
	RequirementOK      RequirementStatus = "ok"
	RequirementProblem RequirementStatus = "problem"
	RequirementUnknown RequirementStatus = "unknown"
)

type Requirement struct {
	Type    string            `json:"type"`
	Status  RequirementStatus `json:"status"`
	Message string            `json:"message"`
	Source  *RuleSource       `json:"source,omitempty"`
}

type CheckResult struct {
	Status       Status        `json:"status"`
	Requirements []Requirement `json:"requirements"`
}
