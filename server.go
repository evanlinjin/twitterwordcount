package main

var chTweets chan string = make(chan string)

func main() {
	// Tweet Getter >>
	tweetGetter := makeTweetGetter("oauth.json", 10000)
	go tweetGetter.start()

	// Tweet Transformer >>
	tweetTransformer := makeTweetTransformer(tweetGetter.chTweets)
	go tweetTransformer.start()

	// Word Storer >>
	wordStorer := makeWordStorer(tweetTransformer.chOut)
	go wordStorer.start()

	// API Server >>
	wordApi := makeWordApi(":8080", wordStorer)
	wordApi.start()
}
