package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"github.com/monory/messenger/database"
)

const (
	selectorSize  = 16
	validatorSize = 32
)

func Register(db *sql.DB, username, password string) error {
	userExists, err := database.CheckUserExists(db, username)
	if err != nil {
		return err
	}
	if userExists {
		return ErrUserExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	err = database.AddUser(db, username, passwordHash)
	if err != nil {
		return err
	}

	return nil
}

func Login(db *sql.DB, username, password string) (*UserToken, error) {
	var t *UserToken
	dbPasswordHash, err := database.GetPasswordHash(db, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, ErrUsernameNotFound
		}
		return t, err
	}

	err = bcrypt.CompareHashAndPassword(dbPasswordHash, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return t, ErrInvalidPassword
		}
		return t, err
	}

	t = NewUserToken()
	for selectorExists := true; selectorExists; {
		t.Random()
		selectorExists, err = database.CheckSelectorExists(db, t.Selector)
		if err != nil {
			return t, err
		}
	}

	dbToken := t.DBToken()
	userID, err := database.GetUserID(db, username)
	if err != nil {
		return t, err
	}
	dbToken.UserID = userID

	err = database.AddToken(db, dbToken)
	if err != nil {
		return t, err
	}

	return t, nil
}

func CheckUserToken(db *sql.DB, t *UserToken) error {
	dbToken, err := database.GetToken(db, t.Selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrInvalidToken
		}
		return err
	}

	hash := sha256.Sum256(t.Validator)
	if subtle.ConstantTimeCompare(hash[:], dbToken.Token) == 1 {
		return nil
	}

	return ErrInvalidToken
}

func MakeChatToken(db *sql.DB, t *UserToken) (*UserToken, error) {
	result := NewUserToken()
	for selectorExists := true; selectorExists; {
		result.Random()
		var err error
		selectorExists, err = database.CheckChatSelectorExists(db, result.Selector)
		if err != nil {
			return result, err
		}
	}

	dbToken, err := database.GetToken(db, t.Selector)
	if err != nil {
		return result, err
	}

	chatDBToken := result.DBToken()
	chatDBToken.UserID = dbToken.UserID
	err = database.AddChatToken(db, chatDBToken)
	if err != nil {
		return result, err
	}

	return result, nil
}

func CheckChatToken(db *sql.DB, t *UserToken) (string, error) {
	result := ""

	dbToken, err := database.GetChatToken(db, t.Selector)
	if err != nil {
		return result, err
	}
	err = database.UseChatToken(db, t.Selector)
	if err != nil {
		return result, err
	}
	result, err = database.GetUsername(db, dbToken.UserID)
	if err != nil {
		return result, err
	}

	return result, nil
}
