package chat

import (
	"log"

	"github.com/monory/messenger/database"
)

func Start() {
	db := database.ConnectDatabase("user=chat_backend password=chat_backend_password dbname=messenger")
	err := db.Ping()
	if err != nil {
		log.Print(err)
	}

	server := NewServer("/ws")
	go server.Listen(db)
}
