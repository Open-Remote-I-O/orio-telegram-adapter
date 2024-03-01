package ports

import (
	"context"
)

type DeviceControlService interface {
	StartServer(ctx context.Context)
}

type RemoteControlService interface {
	StartServer(ctx context.Context)
}
