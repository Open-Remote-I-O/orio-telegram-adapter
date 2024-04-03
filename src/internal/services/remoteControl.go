package services

import (
	"context"
	"orio-telegram-adapter/src/internal/domain/ports"
)

type RemoteControlService struct {
	svc ports.RemoteControlService
}

func NewRemoteControlService(
	remoteControlService ports.RemoteControlService,
) RemoteControlService {
	return RemoteControlService{
		svc: remoteControlService,
	}
}

func (s *RemoteControlService) StartServer(ctx context.Context) {
	s.svc.StartServer(ctx)
}
