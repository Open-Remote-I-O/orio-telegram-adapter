package ports

import (
	"context"
)

type EventEvaluator interface {
	Receive(ctx context.Context)
	Publish(ctx context.Context)
}
