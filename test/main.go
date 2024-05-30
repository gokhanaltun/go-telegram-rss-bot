package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

func main() {
	fp := gofeed.NewParser()

	ticker := time.NewTicker(time.Second * 5)

	for {
		newFeeds := []*gofeed.Item{}

		f, err := fp.ParseURL("https://www.upwork.com/ab/feed/jobs/rss?paging=0-10&payment_verified=1&q=(%22web%20scraping%22)&sort=recency&api_params=1&securityToken=968f3483a6690ece0fdfb61ab5d200c5d322e565e9d1c5cea6f0636c72198c4386c5dbe986f7f2ae1975ffc417b3ac99d8a2aaf3e44986e415f1892129ed2eb8&userUid=1707513677085786112&orgUid=1707513677085786113")
		if err != nil {
			log.Println(err)
		}

		for _, item := range f.Items {
			parsedTime, _ := time.Parse("2006-01-02 15:04:05 +0000 MST", "2024-05-30 10:32:31 +0000 UTC")

			if item.PublishedParsed.After(parsedTime) {
				newFeeds = append(newFeeds, item)
			}
		}

		fmt.Println(len(f.Items), len(newFeeds))
		<-ticker.C
	}

}
