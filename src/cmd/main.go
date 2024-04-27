package main

import (
	"context"
	"os"
	"sync"

	"orio-telegram-adapter/src/internal/adapters"
	"orio-telegram-adapter/src/internal/config"
	"orio-telegram-adapter/src/internal/services"

	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	logger.Debug().
		Msg("logger was configured and instantiated successfully")

	remotedevicecontrollerconf := config.NewDeviceConfig()

	deviceControlAdapter, err := adapters.NewDeviceRemoteControlAdapter(&logger, remotedevicecontrollerconf)
	if err != nil {
		logger.Err(err).Msg("unexpected error while initializing remote device control")
		return
	}

	deviceRemoteService := services.NewDeviceControlService(deviceControlAdapter)

	logger.Debug().
		Msg("Device remote control service configured and instantiated successfully")

	remoteControlAdapter, err := adapters.NewTelegramRemoteControlAdapter(&logger)
	if err != nil {
		logger.Err(err).Msg("unexpected error while initializing telegram connection")
		return
	}

	remoteControlService := services.NewRemoteControlService(
		remoteControlAdapter,
	)

	logger.Debug().
		Msg("Device remote control service configured and instantiated successfully")

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go deviceRemoteService.StartServer(context.Background())

	wg.Add(1)
	go remoteControlService.StartServer(context.Background())

	wg.Wait()
}
