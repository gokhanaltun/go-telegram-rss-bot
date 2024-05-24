package worker

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func parse() {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL("")
	if err != nil {
		panic(err)
	}

	fmt.Println(feed.Title)

}
