package ports

import (
	"context"
)

type RemoteControlService interface {
	StartServer(ctx context.Context)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
