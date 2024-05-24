package migration

import (
	"github.com/gokhanaltun/go-telegram-rss-bot/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	models := []interface{}{
		&models.Feed{},
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			panic("failed to migrate model: " + err.Error())
		}
	}
}
