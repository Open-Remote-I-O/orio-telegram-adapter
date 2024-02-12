package adapters

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	bot "github.com/go-telegram/bot"
	bot_model "github.com/go-telegram/bot/models"
)

type TelegramHandler struct{
	logger *zerolog.Logger
	server *bot.Bot 
}

func NewTelegramRemoteControlAdapter(
	logger *zerolog.Logger,
) (TelegramHandler,error){
	botApiKey, envIsPresent := os.LookupEnv("BOT_API_KEY")
	if !envIsPresent {
		fmt.Println("missing env variable")
	}


	b, err := bot.New(
		botApiKey, 
		[]bot.Option{
			bot.WithDefaultHandler(handler),
		}...,
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	
	return TelegramHandler{
		logger: logger,
		server: b,

	},nil
}

func (th *TelegramHandler) StartServer(ctx context.Context) {
	th.server.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *bot_model.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
	if err != nil{
		log.Error().Err(err).Msg("something went wrong while sending echo to client")
	}
}

