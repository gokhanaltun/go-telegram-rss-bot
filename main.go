package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	bot "github.com/gokhanaltun/go-telegram-bot"
	"github.com/gokhanaltun/go-telegram-rss-bot/commands"
	"github.com/gokhanaltun/go-telegram-rss-bot/middlewares"
	"github.com/gokhanaltun/go-telegram-rss-bot/worker"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithMiddlewares(middlewares.CheckUser),
		bot.WithCallbackQueryDataHandler("deleteRssSelect", bot.MatchTypePrefix, commands.DeleteRssSelectHandler),
		bot.WithCallbackQueryDataHandler("deleteRssConfirm", bot.MatchTypePrefix, commands.DeleteRssConfirmHandler),
	}

	b, err := bot.New(os.Getenv("TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	commands.Init()

	stages := map[int]bot.HandlerFunc{
		commands.RssNameStage: commands.RssNameHandler,
		commands.RssUrlStage:  commands.RssUrlHandler,
	}

	convEnd := bot.ConversationEnd{
		Command:  "/cancel",
		Function: commands.CancelConversation,
	}

	convHandler := bot.NewConversationHandler(stages, &convEnd)

	b.AddConversationHandler(convHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, commands.Start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add", bot.MatchTypeExact, commands.AddRss)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, commands.ListRss)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/delete", bot.MatchTypeExact, commands.DeleteRss)

	go worker.StartRssWorker(rssWorkerCallback)

	b.Start(ctx)
}

func rssWorkerCallback(feeds []*gofeed.Item) {
	for _, item := range feeds {
		fmt.Println(item.PublishedParsed)
	}
}
