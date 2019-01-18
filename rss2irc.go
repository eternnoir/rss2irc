package rss2irc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	irc "github.com/fluffle/goirc/client"
)

func NewRss2Irc(cfg *IrcClientConfig) *Rss2Irc {
	irccfg := irc.NewConfig(cfg.Nick)
	if cfg.Ssl {
		irccfg.SSL = true
		irccfg.SSLConfig = &tls.Config{ServerName: cfg.ServerName}
	}
	irccfg.Server = cfg.Server
	irccfg.NewNick = func(n string) string { return n + "^" }
	conn := irc.Client(irccfg)
	rr := int64(30)
	if cfg.RefreshRate > 0 {
		rr = cfg.RefreshRate
	}
	return &Rss2Irc{
		ircClient:    conn,
		channelFeeds: cfg.GetChannelMap(),
		quit:         make(chan bool),
		RefreshRate:  rr,
	}
}

type Rss2Irc struct {
	ircClient    *irc.Conn
	channelFeeds map[string][]*FeedFetcher
	quit         chan bool
	RefreshRate  int64
}

func (r2i *Rss2Irc) Start(ctx context.Context) error {
	r2i.setHandler()

	fmt.Println("Start connect to server.")
	// Login to irc
	if err := r2i.ircClient.Connect(); err != nil {
		return err
	}

	<-r2i.quit
	return nil
}

func (r2i *Rss2Irc) setHandler() {
	r2i.ircClient.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			fmt.Println("Connected to server.")
			for c := range r2i.channelFeeds {
				conn.Join(c)
				fmt.Printf("[SYS] Joined channel %s.\n", c)
			}
			for k, fetcherList := range r2i.channelFeeds {
				for _, fetcher := range fetcherList {
					fmt.Println(k)
					go r2i.fetchNewFeeds(k, fetcher)
				}
			}
		})
	// And a signal on disconnect
	r2i.ircClient.HandleFunc(irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) { r2i.quit <- true })

	r2i.ircClient.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			fmt.Println("Get new message.")
			fmt.Println(line.Raw)
		})

	r2i.ircClient.HandleFunc(irc.CTCP,
		func(conn *irc.Conn, line *irc.Line) {
			fmt.Printf("[CTCP]%s\n", line.Raw)
		})
}

func (r2i *Rss2Irc) fetchNewFeeds(channelName string, fetcher *FeedFetcher) {
	_, err := fetcher.GetUpdates()
	if err != nil {
		fmt.Printf("[ERROR] %s \n", err.Error())
	}
	for {
		time.Sleep(time.Duration(r2i.RefreshRate) * time.Second)
		updates, err := fetcher.GetUpdates()
		if err != nil {
			fmt.Printf("[ERROR] %s \n", err.Error())
			continue
		}
		for _, u := range updates {
			fmt.Printf("[SEND] %s : %s\n", channelName, u.Title)
			r2i.ircClient.Privmsgf(channelName, "%s | %s", u.Title, u.Link)
		}
	}
}
