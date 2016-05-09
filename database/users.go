package database

import "database/sql"

func AddUser(db *sql.DB, username string, passwordHash []byte) error {
	_, err := db.Exec("INSERT INTO Users (username, password_hash) VALUES ($1, $2)", username, passwordHash)
	return err
}

func CheckUserExists(db *sql.DB, username string) (bool, error) {
	var result bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username).Scan(&result)
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
