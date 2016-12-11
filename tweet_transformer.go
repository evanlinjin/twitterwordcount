package main

import (
	"fmt"
	"strings"
)

func makeTweetTransformer(chTweets chan string) *TweetTransformer {
	transformer := TweetTransformer{
		chIn:   chTweets,
		chOut:  make(chan map[string]int),
		chStop: make(chan int),
	}

	return &transformer
}

type TweetTransformer struct {
	chIn   chan string
	chOut  chan map[string]int
	chStop chan int
}

func (tt *TweetTransformer) start() {
	for {
		select {
		case msg := <-tt.chIn:
			fmt.Println("Msg:", msg)
			tt.chOut <- msgToWordMap(msg)
		case <-tt.chStop:
			return
		}
	}
}

func (tt *TweetTransformer) stop() {
	tt.chStop <- 0
}

// TODO: Fix algorithm.
func msgToWordMap(msg string) map[string]int {
	wordMap := map[string]int{}

	// Trim, lowercase and split msg.
	msg = strings.Trim(strings.ToLower(msg), " ")
	splitMsg := strings.Split(msg, " ")

	// Remove 'Reply to' part.
	if splitMsg[0] == "rt" &&
		strings.HasPrefix(splitMsg[1], "@") &&
		strings.HasSuffix(splitMsg[1], ":") {
		splitMsg = splitMsg[2:]
	}

	// Remove Tags, Hashtags and URLs.
	for i, wrd := range splitMsg {

		if strings.HasPrefix(wrd, "@") ||
			strings.HasPrefix(wrd, "#") ||
			strings.Contains(wrd, "https://") ||
			strings.Contains(wrd, "http://") {
			splitMsg[i] = ""
			continue
		}

		// Deal with words regarding "I"
		isDone := true
		switch wrd {
		case "i":
			splitMsg[i] = "I"
			break
		case "i'll":
			splitMsg[i] = "I'll"
			break
		case "i'm":
			splitMsg[i] = "I'm"
			break
		case "i've":
			splitMsg[i] = "I've"
			break
		default:
			isDone = false
		}

		if isDone {
			continue
		}

		// Deal with Non-Latin and Irrelevant Charactors.
		newWord := []byte{}
		for j := 0; j < len(wrd); j++ {
			if isLatin(wrd[j]) || wrd[j] == '\'' {
				newWord = append(newWord, wrd[j])
			} else {
				newWord = append(newWord, ' ')
			}
		}
		newWordSplit := strings.Split(string(newWord), " ")
		for _, newWrd := range newWordSplit {
			splitMsg = append(splitMsg, newWrd)
		}
	}

	// Fill wordMap >>
	for _, wrd := range splitMsg {
		if wrd != "" && wrd != " " {
			wordMap[wrd] += 1
		}
	}

	//fmt.Println("[msg]:", msg)
	//fmt.Println("[splitMsg]:", splitMsg)
	//fmt.Println("[wordMap]:", wordMap)

	return wordMap
}

func isLatin(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}
