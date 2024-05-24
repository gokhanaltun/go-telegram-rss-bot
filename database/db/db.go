package db

import (
	"github.com/gokhanaltun/go-telegram-rss-bot/database/migration"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDb() *gorm.DB {
	if db == nil {
		gormDB, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		db = gormDB
		migration.Migrate(db)
	}
	return db
}
