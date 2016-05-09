package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math"

	"github.com/monory/messenger/database"
)

type UserToken struct {
	Selector  []byte
	Validator []byte
}

func NewUserToken() *UserToken {
	return &UserToken{
		Selector:  make([]byte, selectorSize),
		Validator: make([]byte, validatorSize),
	}
}

func (t *UserToken) String() string {
	return base64.URLEncoding.EncodeToString(t.Selector) + base64.URLEncoding.EncodeToString(t.Validator)
}

func (t *UserToken) FromString(s string) error {
	selectorEncoded := int(math.Ceil(float64(selectorSize)/3) * 4)
	validatorEncoded := int(math.Ceil(float64(validatorSize)/3) * 4)

	if len(s) != selectorEncoded+validatorEncoded {
		return ErrInvalidToken
	}

	var err error
	t.Selector, err = base64.URLEncoding.DecodeString(s[:selectorEncoded])
	if err != nil {
		return ErrInvalidToken
	}

	t.Validator, err = base64.URLEncoding.DecodeString(s[selectorEncoded:])
	if err != nil {
		return ErrInvalidToken
	}

	return nil
}

func (t *UserToken) Random() {
	rand.Read(t.Selector)
	rand.Read(t.Validator)
}

func (t UserToken) DBToken() *database.DBToken {
	hash := sha256.Sum256(t.Validator)
	return &database.DBToken{
		Selector: t.Selector,
		Token:    hash[:],
	}
}
