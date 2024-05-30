// TODO: The code will be optimized and organized.
package worker

import (
	"log"
	"sort"
	"time"

	"github.com/gokhanaltun/go-telegram-rss-bot/database"
	dbModels "github.com/gokhanaltun/go-telegram-rss-bot/models"
	"github.com/mmcdole/gofeed"
)

func StartRssWorker(callback func(feeds []*gofeed.Item)) {
	db := database.GetDb()
	fp := gofeed.NewParser()

	ticker := time.NewTicker(time.Minute * 1)

	for {
		userFeeds := []dbModels.Feed{}
		result := db.Where(&dbModels.Feed{Notification: true}).Find(&userFeeds)
		if result.Error != nil {
			log.Println("database error", result.Error)
		}

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

		sort.Slice(newFeeds, func(i, j int) bool {
			return newFeeds[i].PublishedParsed.Before(*newFeeds[j].PublishedParsed)
		})

		for _, userFeed := range userFeeds {
			result := db.Where(userFeed).Update("LastRead", newFeeds[len(newFeeds)-1].PublishedParsed.String())
			if result.Error != nil {
				log.Println("database update last read error: ", result.Error)
			}
		}
		callback(newFeeds)
		<-ticker.C
	}
}
