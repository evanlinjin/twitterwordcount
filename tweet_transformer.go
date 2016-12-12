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
			tt.chOut <- msgToWordMap(msg)
		case <-tt.chStop:
			return
		}
	}
}

func (tt *TweetTransformer) stop() {
	tt.chStop <- 0
}

func msgToWordMap(msg string) map[string]int {
	wordMap := map[string]int{}

	// Trim 'msg'.
	msg = strings.Trim(msg, " ")

	// Remove reply-to part.
	if msg[0:4] == "RT @" {
		i := 4
	REPLYTO_LOOP:
		for {
			switch msg[i] {
			case ' ':
				break REPLYTO_LOOP
			case ':':
				msg = msg[i+2 : len(msg)]
				break REPLYTO_LOOP
			}
			i += 1
		}
	}

	// 'wordArr' is an Array of the words split from 'wordMap'.
	// The max n. of words possilbe is n. of chars in msg.
	wordArr, waI := make([]string, len(msg)), 0

	// Loop through msg and extract words.
	// nSpc, nLtn: Counts n. of consec space, latin chars.
	// isTag: Records if current word is a Tag, HashTag.
	// isUrl: Records if current word is an Url.
	nSpc, nLtn, isTag, isUrl := 1, 0, false, false

MSG_LOOP:
	for i := 0; i < len(msg); i++ {
		// Extract current & prev char.
		cc, pc := msg[i], byte(' ')
		if i != 0 {
			pc = msg[i-1]
		}

		// Determine setting of isTag & isUrl.
		if (cc == '@' || cc == '#') && (pc == ' ') {
			isTag = true
		} else if cc == 'h' && pc == ' ' {
			j := i
		URL_LOOP:
			for {
				if j == len(msg) || msg[j] == ' ' {
					// Test URL.
					isUrl = isValidUrl(msg[i:j])
					if j == len(msg) {
						break MSG_LOOP
					} else {
						break URL_LOOP
					}
				}
				j += 1
			}
		} else if cc == ' ' {
			isTag, isUrl = false, false
		}

		// Extract words if not Tag or Url.
		if !isTag && !isUrl {
			if isWordPart(msg[i]) {
				wordArr[waI] += string(msg[i])
				nSpc, nLtn = 0, nLtn+1
			} else {
				if nSpc == 0 {
					waI += 1
				}
				nSpc, nLtn = nSpc+1, 0
			}
		}
	}

	// Add all words of 'wordArr' to 'wordMap'.
	for _, wrd := range wordArr {
		if wrd != "" {
			wordMap[wrd] += 1
		}
	}

	return wordMap
}

func isLatin(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func isNumber(c byte) bool {
	return (c >= '0' && c <= '9')
}

func isWordPart(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c == '\'')
}

func printIVU(in ...interface{}) {
	fmt.Println("[func isValidUrl]", in)
}

func isValidUrl(url string) bool {
	url = strings.TrimSpace(url)
	s, e := [4]int{0, -1, -1, -1}, [4]int{-1, -1, -1, len(url)}
	colonCnt, slashCnt, dotCnt := 0, 0, 0

	for i := 0; i < len(url); i++ {
		c := url[i]
		switch c {
		case ':':
			colonCnt += 1
			if colonCnt > 1 {
				//printIVU("Returning false as: colonCnt > 1")
				return false
			}
			e[0] = i

		case '/':
			slashCnt += 1
			if slashCnt == 2 {
				s[1] = i + 1
			} else if slashCnt == 3 {
				e[2], s[3] = i, i+1
			} else if slashCnt > 3 {
				//printIVU("Returning false as: slashCnt > 3")
				return false
			}

		case '.':
			dotCnt += 1
			if dotCnt == 1 {
				e[1], s[2] = i, i+1
			} else if dotCnt > 1 {
				//printIVU("Returning false as: dotCnt > 1")
				return false
			}

		default:
			if !isLatin(c) && !isNumber(c) {
				//printIVU("Returning false as: !isLatin(c) && !isNumber(c)")
				return false
			}
		}
	}

	for i := 1; i < 4; i++ {
		if s[i] <= s[i-1] || e[i] <= e[i-1] {
			//printIVU("Returning false as: s[i] <= s[i-1] || e[i] <= e[i-1]")
			return false
		}
	}
	return true
}
