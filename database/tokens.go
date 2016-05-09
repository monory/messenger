package database

import "database/sql"

type DBToken struct {
	UserID   int64
	Selector []byte
	Token    []byte
}

func CheckSelectorExists(db *sql.DB, selector []byte) (bool, error) {
	var result bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM tokens WHERE selector=$1)", selector).Scan(&result)
	return result, err
}

func AddToken(db *sql.DB, t *DBToken) error {
	_, err := db.Exec("INSERT INTO tokens (user_id, selector, token) VALUES ($1, $2, $3)", t.UserID, t.Selector, t.Token)
	return err
}

func GetToken(db *sql.DB, selector []byte) (DBToken, error) {
	var result DBToken
	row := db.QueryRow("SELECT user_id, selector, token FROM tokens WHERE selector=$1 AND expires>NOW()", selector)
	err := row.Scan(&result.UserID, &result.Selector, &result.Token)
	if err != nil {
		return result, err
	}

	return result, nil
}
