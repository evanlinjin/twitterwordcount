package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func makeWordApi(port string, wordStorer *WordStorer) *WordApi {
	wordApi := WordApi{
		port:       port,
		wordStorer: wordStorer,
	}
	return &wordApi
}

type WordApi struct {
	port       string
	wordStorer *WordStorer
}

func (wa *WordApi) start() {
	http.HandleFunc("/v1/words", wa.v1WordsEndpoint)
	go log.Fatal(http.ListenAndServe(":8080", nil))
}

func (wa *WordApi) v1WordsEndpoint(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		count = 0
	}
	fmt.Fprintln(w, wa.wordStorer.getWords(count))
}
