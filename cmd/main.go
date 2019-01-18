package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"rss2irc"
	"time"
)

func main() {
	quit := make(chan bool)
	cfgPathPtr := flag.String("c", "./config.json", "Config file path")
	flag.Parse()
	fmt.Printf("[SYS] Load config from %s\n", *cfgPathPtr)
	//var cfgList []*rss2irc.IrcClientConfig

	cfgList, err := getCfgList(*cfgPathPtr)
	if err != nil {
		panic(err)
	}
	for _, cfg := range cfgList {
		r2i := rss2irc.NewRss2Irc(cfg)
		go start(r2i)
	}
	<-quit
}

func getCfgList(path string) ([]*rss2irc.IrcClientConfig, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var cfgList []*rss2irc.IrcClientConfig
	if err := json.Unmarshal([]byte(byteValue), &cfgList); err != nil {
		return nil, err
	}
	return cfgList, nil
}

func start(r2i *rss2irc.Rss2Irc) {
	for {
		if err := r2i.Start(context.Background()); err != nil {
			fmt.Printf("[ERROR] Start error %s\n", err)
		}
		time.Sleep(30 * time.Second)
	}
}
