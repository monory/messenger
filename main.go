package main

import (
	"log"
	"net/http"

	"github.com/monory/messager-backend/chat"
	"github.com/monory/messager-backend/web"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := chat.NewServer("/entry")
	go server.Listen()

	// login-register server
	http.HandleFunc("/register", web.Handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
