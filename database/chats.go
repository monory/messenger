package database

import (
	"database/sql"
	"log"
	"time"
)

type SimpleMessage struct {
	Author    string `json:"author"`
	Message   string `json:"message"`
	Timestamp time.Time
}

func GetContacts(db *sql.DB, id int64) []string {
	var result []string

	rows, err := db.Query("SELECT DISTINCT u.username FROM contacts c JOIN users u ON c.user_id=u.id WHERE c.owner_id = $1", id)
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

func GetMessages(db *sql.DB, id int64, contact string) []SimpleMessage {
	var result []SimpleMessage

	rows, err := db.Query("SELECT u.username, m.message, m.time_stamp FROM messages m JOIN users u ON m.author_id=u.id WHERE author_id=$1 AND receiver_id=(SELECT id FROM users WHERE username=$2) OR author_id=(SELECT id FROM users WHERE username=$2) AND receiver_id=$1 ORDER BY m.time_stamp", id, contact)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var t SimpleMessage
		err = rows.Scan(&t.Author, &t.Message, &t.Timestamp)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, t)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func AddMessage(db *sql.DB, msg *TextMessage) {
	_, err := db.Exec(
		"INSERT INTO messages (message, author_id, receiver_id) VALUES ($1, (SELECT id FROM users WHERE username=$2), (SELECT id FROM users WHERE username=$3))", msg.Message, msg.Author, msg.Receiver)
	if err != nil {
		log.Println(err)
	}
}

func AddContact(db *sql.DB, id int64, contact string) error {
	_, err := db.Exec("INSERT INTO contacts (owner_id, user_id) VALUES ($1, (SELECT id FROM users WHERE username=$2))", id, contact)
	return err
}
