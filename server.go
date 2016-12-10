package main

import (
	"fmt"
	"log"
	"net/http"
)

var chTweets chan string = make(chan string)

func main() {
	// Tweet Getter >>
	tweetGetter := makeTweetGetter("oauth.json", 10000)
	go tweetGetter.start()

	// Tweet Transformer >>
	tweetTransformer := makeTweetTransformer(tweetGetter.chTweets)
	go tweetTransformer.start()

	// API Server >>
	http.HandleFunc("/v1/words", v1Words)
	go log.Fatal(http.ListenAndServe(":8080", nil))

}

func v1Words(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[REQUEST]", r.URL.String())
	msg := <-chTweets
	fmt.Println("[RESPONSE]", msg)
	fmt.Fprintln(w, msg)
}
