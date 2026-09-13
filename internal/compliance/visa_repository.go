package compliance

import (
	"context"
	"time"
)

type VisaRuleRepository interface {
	FindVisaRule(
		ctx context.Context,
		nationalityCode string,
		destinationCode string,
		departureDate time.Time,
	) (VisaRule, error)
}
