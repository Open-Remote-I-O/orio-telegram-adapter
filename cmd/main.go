package main

import (
	"context"
	"os"
	"os/signal"

	"orio-telegram-adapter/internal/adapters"
	"orio-telegram-adapter/internal/config"
	"orio-telegram-adapter/internal/services"

	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

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
		Msg("Device remote service configured correctly")

	// remoteControlAdapter, err := adapters.NewTelegramRemoteControlAdapter(&logger)
	// if err != nil {
	// 	logger.Err(err).Msg("unexpected error while initializing telegram connection")
	// 	return
	// }

	// remoteControlService := services.NewRemoteControlService(
	// 	remoteControlAdapter,
	// )

	// logger.Debug().
	// 	Msg("Remove control server configured correctly")

	go deviceRemoteService.Start(ctx)

	// go remoteControlService.StartServer(ctx)

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()

	logger.Fatal().Msg("Interrupting")

	os.Exit(0)
	// log.Println("shutting down gracefully, press Ctrl+C again to force")

	// // Perform application shutdown with a maximum timeout of 5 seconds.
	// timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// go func() {
	// 	if err := server.Shutdown(timeoutCtx); err != nil {
	// 		logger.Fatal().Msg(err)
	// 	}
	// }()

	// select {
	// case <-timeoutCtx.Done():
	// 	if timeoutCtx.Err() == context.DeadlineExceeded {
	// 		log.Fatalln("timeout exceeded, forcing shutdown")
	// 	}

	// 	os.Exit(0)
	// }

	// stop()
}
