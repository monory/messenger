package main

import (
	"log"
	"net/http"

	"github.com/monory/messager-backend/chat"
	"github.com/monory/messager-backend/database"
	"github.com/monory/messager-backend/web"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := chat.NewServer("/entry")
	go server.Listen()

	// login-register server
	db := database.ConnectDatabase("???:???@/messenger")
	http.HandleFunc("/register", web.MakeHandler(web.Handler, db))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
