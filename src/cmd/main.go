package main

import (
	"context"
	"fmt"
	"orio-telegram-adapter/src/internal/adapters"
	"orio-telegram-adapter/src/internal/services"
	"os"

	"github.com/rs/zerolog"
)

// Send any text message to the bot after the bot has been started

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	logger.Debug().
		Msg("logger was configured and instantiated successfully")

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
 
	remoteControlService.StartServer(context.TODO())

	logger.Debug().
		Msg("remote control service started")
}
