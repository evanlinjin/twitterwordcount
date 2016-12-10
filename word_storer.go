package main

import "fmt"

func makeWordStorer(chWrdCnts chan map[string]int) WordStorer {
	wordStorer := WordStorer{
		dict: make(map[string]int),
		chIn: chWrdCnts,
	}
	return wordStorer
}

type WordStorer struct {
	dict   map[string]int
	chIn   chan map[string]int
	chStop chan int
}

func (ws *WordStorer) start() {
	for {
		select {
		case wrdCnts := <-ws.chIn:
			for wrd, cnt := range wrdCnts {
				ws.dict[wrd] += cnt
			}
			fmt.Println(ws.dict)
		case <-ws.chStop:
			return
		}
	}
}

func (ws *WordStorer) stop() {
	ws.chStop <- 0
}
