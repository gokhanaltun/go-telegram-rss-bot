// TODO: "Command will be added to set whether notification messages will be received."
package commands

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	bot "github.com/gokhanaltun/go-telegram-bot"
	"github.com/gokhanaltun/go-telegram-bot/models"
	"github.com/gokhanaltun/go-telegram-rss-bot/database"
	dbModels "github.com/gokhanaltun/go-telegram-rss-bot/models"
	"github.com/mmcdole/gofeed"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	rssName string
)

const (
	RssNameStage = iota
	RssUrlStage
)

func Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	startMessage := "Merhaba " + update.Message.From.FirstName + " :)\n\n" +
		"Yeni kayıt eklemek için      /add\n\n" +
		"Kayıtları listelemek için    /list\n\n" +
		"Bir kayıt silmek için          /delete"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   startMessage,
	})
}

func AddRss(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SetActiveConversationStage(RssNameStage)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Rss ismi girin. İptal etmek için /cancel komutunu kullanın.",
	})
}

func RssNameHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	result := db.Where("Name = ?", update.Message.Text).First(&dbModels.Feed{})
	if result.Error != gorm.ErrRecordNotFound {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Bu isimde bir kayıt zaten var.",
		})
		return
	}

	rssName = update.Message.Text

	b.SetActiveConversationStage(RssUrlStage)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Rss url girin.",
	})
}

func RssUrlHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fp := gofeed.NewParser()

	_, err := fp.ParseURL(update.Message.Text)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Rss eklenirken bir hata oluştu.\n" + fmt.Sprint(err),
		})
		return
	}

	result := db.Where("Url = ?", update.Message.Text).First(&dbModels.Feed{})
	if result.Error != gorm.ErrRecordNotFound {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Bu url'e sahip bir kayıt zaten var.",
		})
		return
	}

	result = db.Create(&dbModels.Feed{
		Name:         rssName,
		Url:          update.Message.Text,
		LastRead:     "2006-01-02 15:04:05 +0000 UTC",
		Notification: true,
	})

	if result.Error != nil {
		b.EndConversation()
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Rss eklenemedi.\n" + fmt.Sprint(result.Error),
		})
		return
	}

	b.EndConversation()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Rss eklendi.",
	})
}

func ListRss(ctx context.Context, b *bot.Bot, update *models.Update) {
	feeds := []dbModels.Feed{}
	result := db.Find(&feeds)
	if result.Error != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Bir hata oluştu.\n" + result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Hiç rss kaydı yok.",
		})
		return
	}

	for _, feed := range feeds {
		name := "Name: " + feed.Name + "\n"
		url := "Url: " + feed.Url + "\n"
		lastRead := "Last Read: " + feed.LastRead + "\n"
		notification := "Notification: " + map[bool]string{true: "true", false: "false"}[feed.Notification]

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   name + url + lastRead + notification,
		})
	}
}

func DeleteRss(ctx context.Context, b *bot.Bot, update *models.Update) {
	feeds := []dbModels.Feed{}
	result := db.Find(&feeds)
	if result.RowsAffected == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Hiç rss kaydı yok.",
		})
		return
	}

	inlineKeyboard := &models.InlineKeyboardMarkup{}
	inlineKeyboardButtons := [][]models.InlineKeyboardButton{}

	for _, feed := range feeds {
		button := []models.InlineKeyboardButton{
			{
				Text:         feed.Name,
				CallbackData: "deleteRssSelect_" + fmt.Sprint(feed.ID) + "_" + feed.Name,
			},
		}

		inlineKeyboardButtons = append(inlineKeyboardButtons, button)
	}

	inlineKeyboardButtons = append(inlineKeyboardButtons, []models.InlineKeyboardButton{
		{Text: "İptal", CallbackData: "deleteRssSelect_cancel"},
	})

	inlineKeyboard.InlineKeyboard = inlineKeyboardButtons

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Silmek istediğiniz kaydı seçin.",
		ReplyMarkup: inlineKeyboard,
	})

}

func DeleteRssSelectHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:          update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:       update.CallbackQuery.Message.Message.ID,
		InlineMessageID: update.CallbackQuery.InlineMessageID,
		ReplyMarkup:     models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{}},
	})

	data := strings.Split(update.CallbackQuery.Data, "_")
	if data[1] == "cancel" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Silme işlemi iptal edildi.",
		})
		return
	}

	inlineKeyboardMarkup := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Evet", CallbackData: "deleteRssConfirm_" + data[1] + "_" + data[2]},
				{Text: "Hayır", CallbackData: "deleteRssConfirm_0"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        data[2] + " isimli kaydı silmek istiyor musunuz?",
		ReplyMarkup: inlineKeyboardMarkup,
	})
}

func DeleteRssConfirmHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:          update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:       update.CallbackQuery.Message.Message.ID,
		InlineMessageID: update.CallbackQuery.InlineMessageID,
		ReplyMarkup:     models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{}},
	})

	data := strings.Split(update.CallbackQuery.Data, "_")
	if data[1] == "0" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Silme işlemi iptal edildi.",
		})
		return
	}

	feed := dbModels.Feed{}
	intData, _ := strconv.Atoi(data[1])
	feed.ID = uint(intData)

	result := db.Delete(&feed)
	if result.Error != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Silme işlemi başarısız oldu.\n" + fmt.Sprint(result.Error),
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   data[2] + " silindi.",
	})
}

func CancelConversation(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "İptal edildi.",
	})
}

func Init() {
	db = database.GetDb()
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}
