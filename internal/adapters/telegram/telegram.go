package adapters

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"

	"orio-telegram-adapter/internal/domain"

	bot "github.com/go-telegram/bot"
	// bot_model "github.com/go-telegram/bot/models"
)

type TelegramHandler struct {
	logger     *zerolog.Logger
	server     *bot.Bot
	messageBus domain.MessageBus
}

func NewTelegramRemoteControlAdapter(
	logger *zerolog.Logger,
	messageBus domain.MessageBus,
) (*TelegramHandler, error) {
	// botApiKey, envIsPresent := os.LookupEnv("BOT_API_KEY")
	// if !envIsPresent {
	// 	return nil, errors.New("BOT_API_KEY env was not provided")
	// }

	// b, err := bot.New(
	// 	botApiKey,
	// 	[]bot.Option{
	// 		bot.WithDefaultHandler(handler),
	// 	}...,
	// )
	// if err != nil {
	// 	logger.Err(err).Msg("failed to initialize new telegram bot")
	// 	return nil, err
	// }

	// return &TelegramHandler{
	// 	logger: logger,
	// 	server: b,
	// }, nil
	return &TelegramHandler{
		logger:     logger,
		server:     nil,
		messageBus: messageBus,
	}, nil
}

func (th *TelegramHandler) Start(ctx context.Context) {
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			th.logger.
				Warn().
				Int("chan length", len(th.messageBus.CommandCh)).
				Msg("sending mock data message bus channel for consumption")
			th.messageBus.CommandCh <- domain.Command{
				DeviceID:  "6969",
				CommandID: 1,
				// fill with rgb value
				Payload: []byte{},
			}
		}
	}
	// th.logger.Info().Msg("starting remote device control service")
	// th.server.Start(ctx)
}

func (th *TelegramHandler) Stop(ctx context.Context) error {
	panic("nothing")
}

// func handler(ctx context.Context, b *bot.Bot, update *bot_model.Update) {
// 	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID: update.Message.Chat.ID,
// 		Text:   update.Message.Text,
// 	})
// 	if err != nil {
// 		log.Error().Err(err).Msg("something went wrong while sending message to client")
// 	}
// }
