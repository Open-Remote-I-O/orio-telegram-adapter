package ports

import (
	"context"
)

type RemoteControlService interface {
	Start(ctx context.Context)
	Stop(ctx context.Context) error
}
