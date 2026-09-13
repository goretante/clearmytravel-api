package compliance

import (
	"context"
	"time"
)

type StayRuleRepository interface {
	FindStayRule(
		ctx context.Context,
		nationalityCode string,
		destinationCode string,
		departureDate time.Time,
	) (StayRule, error)
}
