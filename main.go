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
	db := database.ConnectDatabase("user=messenger_user password=example_password dbname=messenger")
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// websocket server
	server := chat.NewServer("/entryhandler")
	go server.Listen(db)

	// login-register server
	http.HandleFunc("/registerhandler", web.MakeHandler(web.Handler, db))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
