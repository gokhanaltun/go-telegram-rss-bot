package commands

import (
	"context"

	bot "github.com/gokhanaltun/go-telegram-bot"
	"github.com/gokhanaltun/go-telegram-bot/models"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "default handler",
	})
}

func Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "start handler",
	})
}
