package main

import (
	"context"
	"fmt"
	"orio-telegram-adapter/src/internal/adapters"
	"orio-telegram-adapter/src/internal/services"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	logger.Debug().
		Msg("logger was configured and instantiated successfully")

	deviceControlAdapter, err := adapters.NewDeviceRemoteControlAdapter(&logger)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	deviceRemoteController := services.NewDeviceControlService(&deviceControlAdapter)

	deviceRemoteController.StartServer(context.Background())

	logger.Debug().
		Msg("Device remove control service configured and instantiated successfully")

	remoteControlAdapter, err := adapters.NewTelegramRemoteControlAdapter(&logger)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	remoteControlService := services.NewRemoteControlService(
		&remoteControlAdapter,
	)

	logger.Debug().
		Msg("Remote control service configured and instantiated successfully")

	remoteControlService.StartServer(context.Background())

	logger.Debug().
		Msg("remote control service started")

}
