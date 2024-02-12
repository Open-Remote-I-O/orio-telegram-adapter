package ports

import (
	"context"
)


type RemoteControlService interface {
	StartServer(ctx context.Context)
}
