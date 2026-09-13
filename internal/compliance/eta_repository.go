package compliance

import (
	"context"
	"time"
)

type ETARuleRepository interface {
	FindETARule(
		ctx context.Context,
		nationalityCode string,
		destinationCode string,
		departureDate time.Time,
	) (ETARule, error)
}
