package worker

import (
	"log"

	"github.com/mmcdole/gofeed"
)

func StartRssWorker(callback func(feed *gofeed.Feed)) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL("")
	if err != nil {
		log.Println(err)
	}

	callback(feed)

}
