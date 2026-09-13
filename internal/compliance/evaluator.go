package compliance

import "context"

type Evaluator interface {
	Evaluate(ctx context.Context, req CheckRequest) Requirement
}
