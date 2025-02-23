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

func (s *DeviceControlService) Start(ctx context.Context) {
	err := s.svc.Start(ctx)
	if err != nil {
		panic(err)
	}
}

func (s *DeviceControlService) Stop(ctx context.Context) {
	err := s.svc.Stop(ctx)
	if err != nil {
		panic(err)
	}
}
