package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // to initialize mysql

	"log"

	"golang.org/x/crypto/bcrypt"
)

func ConnectDatabase(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func AddUser(db *sql.DB, username, password string) bool {
	usernameExists := true
	var id uint64
	err := db.QueryRow("SELECT id FROM Users WHERE login=?", username).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			usernameExists = false
		} else {
			log.Fatal(err)
		}
	}

	if usernameExists {
		return false
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10) // 10 выбрано как хорошее умолчание для стоимости
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO Users (login, password_hash) VALUES (?, ?)", username, passwordHash)
	log.Println("Good register:", username, string(passwordHash))
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func CheckUser(db *sql.DB, username, password string) bool {
	var passwordHash []byte
	err := db.QueryRow("SELECT password_hash FROM Users WHERE login=?", username).Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	if err != nil {
		return false
	}

	log.Println("Log in:", username, string(passwordHash))
	return true
}
