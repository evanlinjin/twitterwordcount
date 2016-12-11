package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func makeTweetGetter(path string, maxStore int) *TweetGetter {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	var v interface{}
	json.Unmarshal(fileData, &v)
	vMap := v.(map[string]interface{})

	getter := TweetGetter{
		consumerKey:    vMap["consumer_key"].(string),
		consumerSecret: vMap["consumer_secret"].(string),
		accessToken:    vMap["access_token"].(string),
		accessSecret:   vMap["access_secret"].(string),
		chTweets:       make(chan string, maxStore),
		chStop:         make(chan int, 1),
	}

	return &getter
}

type TweetGetter struct {
	consumerKey, consumerSecret, accessToken, accessSecret string
	chTweets                                               chan string
	chStop                                                 chan int
}

func (tg *TweetGetter) start() {
	if tg.consumerKey == "" || tg.consumerSecret == "" || tg.accessToken == "" || tg.accessSecret == "" {
		log.Fatalln("Consumer key/secret and Access token/secret required.")
	}

	config := oauth1.NewConfig(tg.consumerKey, tg.consumerSecret)
	token := oauth1.NewToken(tg.accessToken, tg.accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client.
	client := twitter.NewClient(httpClient)

	params := &twitter.StreamSampleParams{StallWarnings: twitter.Bool(true)}
	stream, err := client.Streams.Sample(params)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case msg := <-stream.Messages:
			switch msg.(type) {
			case *twitter.Tweet:
				tweet := msg.(*twitter.Tweet)
				if tweet.Lang == "en" {
					tg.chTweets <- tweet.Text
				}
			}
		case <-tg.chStop:
			return
		}
	}
}

func (tg *TweetGetter) stop() {
	tg.chStop <- 0
}
