package services

import (
	"context"
	"orio-telegram-adapter/src/internal/domain/ports"
)

type DeviceControlService struct {
	svc ports.DeviceControlService
}

func NewDeviceControlService(
	deviceControlService ports.DeviceControlService,
) DeviceControlService {
	return DeviceControlService{
		svc: deviceControlService,
	}
}

func (s *DeviceControlService) StartServer(ctx context.Context) {
	s.svc.StartServer(ctx)
}
