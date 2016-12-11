package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

func makeWordStorer(chWrdCnts chan map[string]int) *WordStorer {
	wordStorer := WordStorer{
		dict:   make(map[string]int),
		chIn:   chWrdCnts,
		chStop: make(chan int),
	}
	return &wordStorer
}

type WordStorer struct {
	dict    map[string]int
	muxDict sync.Mutex

	chIn   chan map[string]int
	chStop chan int
}

func (ws *WordStorer) start() {
	for {
		// Store words in Unsorted Dictionary.
		select {
		case wrdCnts := <-ws.chIn:
			ws.muxDict.Lock()
			for wrd, cnt := range wrdCnts {
				ws.dict[wrd] += cnt
			}
			ws.muxDict.Unlock()
			fmt.Println(ws.dict)
		case <-ws.chStop:
			return
		}
	}
}

func (ws *WordStorer) stop() {
	ws.chStop <- 0
}

func (ws *WordStorer) getWords(maxN int) string {
	// Add all elements of 'dict' to 'sortedDict'.
	ws.muxDict.Lock()
	sortedDict := ByCount{}
	for k, v := range ws.dict {
		sortedDict = append(sortedDict, SortableWord{k, v})
	}
	ws.muxDict.Unlock()
	// Sort 'sortedDict'.
	sort.Sort(sort.Reverse(sortedDict))
	// Slice 'sortedDict'.
	sortedDict = sortedDict[:maxN]
	// Turn to json.
	b, err := json.Marshal(sortedDict)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(sortedDict)
	return string(b)
}

// TYPES USED FOR SORTING >>>

type SortableWord struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type ByCount []SortableWord

func (w ByCount) Len() int           { return len(w) }
func (w ByCount) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w ByCount) Less(i, j int) bool { return w[i].Count < w[j].Count }
