package database

import (
	"database/sql"
	"encoding/base64"

	_ "github.com/lib/pq" // to initialize postgres

	"log"

	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

func ConnectDatabase(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func AddUser(db *sql.DB, username, password string) bool {
	usernameExists := true
	var id uint64
	err := db.QueryRow("SELECT id FROM Users WHERE login=$1", username).Scan(&id)
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

	_, err = db.Exec("INSERT INTO Users (login, password_hash) VALUES ($1, $2)", username, passwordHash)
	log.Println("Good register:", username, string(passwordHash))
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func CheckUser(db *sql.DB, username, password string) bool {
	var passwordHash []byte
	err := db.QueryRow("SELECT password_hash FROM Users WHERE login=$1", username).Scan(&passwordHash)
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

func GenerateToken(db *sql.DB, username string) string {
	token := make([]byte, 64)
	rand.Read(token)

	_, err := db.Exec("INSERT INTO Tokens (user_id, token) VALUES ((SELECT id FROM Users WHERE login=$1), $2)",
		username, token)
	if err != nil {
		log.Fatal(err)
	}
	return base64.RawURLEncoding.EncodeToString(token)
}

func CheckToken(db *sql.DB, encodedToken string) (string, error) {
	token, err := base64.RawURLEncoding.DecodeString(encodedToken)
	if err != nil {
		return "", err
	}

	var name string
	err = db.QueryRow("SELECT login FROM Users WHERE Users.id=(SELECT user_id FROM Tokens WHERE token=$1)", token).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		log.Fatal(err)
	}

	return name, nil
}
