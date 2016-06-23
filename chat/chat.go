package chat

import (
	"log"

	"github.com/monory/messenger/database"
)

func Start() {
	db := database.ConnectDatabase("user=chat_backend dbname=messenger")
	err := db.Ping()
	if err != nil {
		log.Print(err)
	}

	server := NewServer("/ws")
	go server.Listen(db)
}
