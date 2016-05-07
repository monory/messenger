package site

import (
	"log"

	"github.com/monory/messenger/database"
)

func Start() {
	db := database.ConnectDatabase("user=web_backend password=web_backend_password dbname=messenger")
	err := db.Ping()
	if err != nil {
		log.Print(err)
	}

	Listen(db)
}
