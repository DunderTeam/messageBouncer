package main

import (
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"time"
)

var Cache = cache.New(5*time.Minute, 5*time.Minute)

func main() {

	http.HandleFunc("/message", HandleMessage)

	err := http.ListenAndServe(":$PORT", nil)
	if err != nil {
		log.Fatal("Server not running", err)
		return
	}
}
