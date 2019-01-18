package rss2irc

type IrcClientConfig struct {
	Server      string           `json:"host"`
	Ssl         bool             `json:"ssl"`
	ServerName  string           `json:"serverName"`
	Nick        string           `json:"nick"`
	RefreshRate int64            `json:"refreshRate"`
	Channels    []*ChannelConfig `json:"channels"`
	Format      string           `json:format`
}

func (icc *IrcClientConfig) GetChannelMap() map[string][]*FeedFetcher {
	ret := make(map[string][]*FeedFetcher)
	for _, cc := range icc.Channels {
		ret[cc.ChannelName] = getFetchers(cc.RssFeeds)
	}
	return ret
}

func getFetchers(feedurls []string) []*FeedFetcher {
	ret := []*FeedFetcher{}
	for _, feed := range feedurls {
		ret = append(ret, NewFeedFetcher(feed))
	}
	return ret
}

type ChannelConfig struct {
	ChannelName string   `json:"channelName"`
	RssFeeds    []string `json:rssFeeds`
}
