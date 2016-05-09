package database

import "database/sql"

func CheckChatSelectorExists(db *sql.DB, selector []byte) (bool, error) {
	var result bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM chat_tokens WHERE selector=$1)", selector).Scan(&result)
	return result, err
}

func AddChatToken(db *sql.DB, t *DBToken) error {
	_, err := db.Exec("INSERT INTO chat_tokens (user_id, selector, token) VALUES ($1, $2, $3)", t.UserID, t.Selector, t.Token)
	return err
}

func GetChatToken(db *sql.DB, selector []byte) (DBToken, error) {
	var result DBToken
	row := db.QueryRow("SELECT user_id, selector, token FROM chat_tokens WHERE selector=$1 AND expires>NOW()", selector)
	err := row.Scan(&result.UserID, &result.Selector, &result.Token)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UseChatToken(db *sql.DB, selector []byte) error {
	_, err := db.Exec("DELETE FROM chat_tokens WHERE selector=$1", selector)
	return err
}
