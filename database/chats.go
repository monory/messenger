package database

import (
	"database/sql"
	"log"
)

func GetChats(db *sql.DB) []string {
	var result []string

	rows, err := db.Query("SELECT name FROM chats")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}
