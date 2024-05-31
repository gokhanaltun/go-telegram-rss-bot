package worker

import (
	"log"
	"sort"
	"time"

	"github.com/gokhanaltun/go-telegram-rss-bot/database"
	dbModels "github.com/gokhanaltun/go-telegram-rss-bot/models"
	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"
)

func StartRssWorker(callback func(feeds []*gofeed.Item)) {
	db := database.GetDb()
	fp := gofeed.NewParser()

	ticker := time.NewTicker(time.Minute * 5)

	for {

		userFeeds, result := getUserFeeds(db)
		if result.Error != nil {
			log.Println("database error (getUserFeeds:rssworker.go): ", result.Error)
		}

		newFeeds := parseFeeds(userFeeds, fp)

		if len(newFeeds) > 0 {
			sortFeedsByPublishedDate(newFeeds)
			lastRead := newFeeds[len(newFeeds)-1].PublishedParsed.String()

			updateFeedLastDates(userFeeds, lastRead, db)

			callback(newFeeds)
		}

		<-ticker.C
	}
}

func getUserFeeds(db *gorm.DB) (userFeeds []dbModels.Feed, dbResult *gorm.DB) {
	feeds := []dbModels.Feed{}
	result := db.Where(&dbModels.Feed{Notification: true}).Find(&feeds)

	return feeds, result
}

func parseFeeds(userFeeds []dbModels.Feed, fp *gofeed.Parser) (feedItems []*gofeed.Item) {
	newFeeds := []*gofeed.Item{}

	for _, userFeed := range userFeeds {
		f, err := fp.ParseURL(userFeed.Url)
		if err != nil {
			log.Println("feed parse error: ", err)
		}

		for _, item := range f.Items {
			parsedTime, err := time.Parse("2006-01-02 15:04:05 +0000 MST", userFeed.LastRead)
			if err != nil {
				log.Println("time parse error: ", err)
			}
			if item.PublishedParsed.After(parsedTime) {
				newFeeds = append(newFeeds, item)
			}
		}
	}

	return newFeeds
}

func sortFeedsByPublishedDate(feeds []*gofeed.Item) {
	sort.Slice(feeds, func(i, j int) bool {
		return feeds[i].PublishedParsed.Before(*feeds[j].PublishedParsed)
	})
}

func updateFeedLastDates(userFeeds []dbModels.Feed, lastRead string, db *gorm.DB) {
	for _, userFeed := range userFeeds {
		result := db.Model(&userFeed).Update("LastRead", lastRead)

		if result.Error != nil {
			log.Println("database update last read error: ", result.Error)
		}
	}
}
