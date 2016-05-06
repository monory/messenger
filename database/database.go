package database

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq" // to initialize postgres
)

func ConnectDatabase(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

type DBToken struct {
	UserID   int64
	Selector []byte
	Token    []byte
	expires  time.Time
}

func AddUser(db *sql.DB, username string, passwordHash []byte) error {
	_, err := db.Exec("INSERT INTO Users (username, password_hash) VALUES ($1, $2)", username, passwordHash)
	return err
}

func CheckUserExists(db *sql.DB, username string) (bool, error) {
	var result bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username).Scan(&result)
	return result, err
}

func CheckSelectorExists(db *sql.DB, selector []byte) (bool, error) {
	var result bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM tokens WHERE selector=$1)", selector).Scan(&result)
	return result, err
}

func GetPasswordHash(db *sql.DB, username string) ([]byte, error) {
	var result []byte
	err := db.QueryRow("SELECT password_hash FROM users WHERE username=$1", username).Scan(&result)
	return result, err
}

func GetUserID(db *sql.DB, username string) (int64, error) {
	var result int64
	err := db.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&result)
	return result, err
}

func AddToken(db *sql.DB, t DBToken) error {
	_, err := db.Exec("INSERT INTO tokens (user_id, selector, token) VALUES ($1, $2, $3)", t.UserID, t.Selector, t.Token)
	return err
}

func GetToken(db *sql.DB, selector []byte) (DBToken, error) {
	var result DBToken
	row := db.QueryRow("SELECT user_id, selector, token, expires FROM tokens WHERE selector=$1 AND expires>NOW()", selector)
	err := row.Scan(&result.UserID, &result.Selector, &result.Token, &result.expires)
	if err != nil {
		return result, err
	}

	return result, nil
}

func validateTokenExpiration(db *sql.DB, t DBToken) error {
	if time.Now().After(t.expires) {
		_, err := db.Exec("DELETE FROM tokens WHERE selector=$1", t.Selector)
		if err != nil {
			return err
		}
		return errors.New("token expired")
	}
	return nil
}
