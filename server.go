package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/v1/words", v1Words)
	go log.Fatal(http.ListenAndServe(":8080", nil))

}

func v1Words(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", r.URL.EscapedPath())
}
