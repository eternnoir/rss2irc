package rss2irc

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func NewFeedFetcher(url string) *FeedFetcher {
	parser := gofeed.NewParser()
	return &FeedFetcher{
		FeedUrl: url,
		parser:  parser,
		tmpMap:  make(map[string]bool),
	}
}

type FeedFetcher struct {
	FeedUrl string
	parser  *gofeed.Parser
	tmpMap  map[string]bool
}

func (ff *FeedFetcher) GetUpdates() ([]*gofeed.Item, error) {
	feed, err := ff.parser.ParseURL(ff.FeedUrl)
	if err != nil {
		return nil, err
	}

	ret := []*gofeed.Item{}
	newTmpMap := make(map[string]bool)
	for _, item := range feed.Items {
		itemKey := fmt.Sprintf("%s;%s;%s", item.Title, item.Description, item.GUID)
		newTmpMap[itemKey] = true
		if _, ok := ff.tmpMap[itemKey]; !ok {
			ret = append(ret, item)
		}
	}
	ff.tmpMap = newTmpMap
	return ret, nil
}
