package auth

import (
	"crypto/rand"
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

func Login(db *sql.DB, username, password string) (UserToken, error) {
	var t UserToken
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

	t, err = generateUserToken(db)
	if err != nil {
		return t, err
	}

	dbToken, err := generateDBToken(db, t, username)
	if err != nil {
		return t, err
	}

	err = database.AddToken(db, dbToken)
	if err != nil {
		return t, err
	}

	return t, nil
}

func generateUserToken(db *sql.DB) (UserToken, error) {
	var result UserToken

	result.Selector = make([]byte, selectorSize)
	for selectorExists := true; selectorExists; {
		rand.Read(result.Selector)
		var err error
		selectorExists, err = database.CheckSelectorExists(db, result.Selector)
		if err != nil {
			return result, err
		}
	}

	result.Validator = make([]byte, validatorSize)
	rand.Read(result.Validator)

	return result, nil
}

func generateDBToken(db *sql.DB, t UserToken, username string) (database.DBToken, error) {
	var result database.DBToken
	var err error

	result.UserID, err = database.GetUserID(db, username)
	if err != nil {
		return result, err
	}

	result.Selector = t.Selector
	hash := sha256.Sum256(t.Validator)
	result.Token = hash[:]
	return result, nil
}

func CheckUserToken(db *sql.DB, t UserToken) error {
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
