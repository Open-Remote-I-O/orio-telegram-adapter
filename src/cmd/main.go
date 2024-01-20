package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"database/sql"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Send any text message to the bot after the bot has been started

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Print("hello world")
	fmt.Println("starting bot")
	db, err := sql.Open("sqlite3", "./docker-volume-config/sql/database.db")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	botApiKey, envIsPresent := os.LookupEnv("BOT_API_KEY")
	if !envIsPresent {
		fmt.Println("missing env variable")
	}

	b, err := bot.New(botApiKey, opts...)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
