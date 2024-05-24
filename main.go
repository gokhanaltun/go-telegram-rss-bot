package main

import (
	"context"
	"os"
	"os/signal"

	bot "github.com/gokhanaltun/go-telegram-bot"
	"github.com/gokhanaltun/go-telegram-rss-bot/commands"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(commands.DefaultHandler),
	}

	b, err := bot.New(os.Getenv("TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, commands.Start)

	b.Start(ctx)
}
