package main

import (
	"log"
	"net/http"

	"github.com/monory/messenger/chat"
	"github.com/monory/messenger/site"
)

func main() {
	site.Start()
	chat.Start()

	log.Println("Start success")
	http.ListenAndServe(":8080", nil)
}
