package ports

import (
	"context"
)

type DeviceControlService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
