package middlewares

import (
	"context"
	"fmt"
	"os"

	bot "github.com/gokhanaltun/go-telegram-bot"
	"github.com/gokhanaltun/go-telegram-bot/models"
)

func CheckUser(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			user_id := fmt.Sprintf("%d", update.Message.From.ID)
			if user_id != os.Getenv("ID") {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "yetki yok",
				})
				return
			}
		}
		next(ctx, b, update)
	}
}
