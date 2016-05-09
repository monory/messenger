package auth

import (
	"encoding/base64"
	"math"
)

type UserToken struct {
	Selector  []byte
	Validator []byte
}

func (t UserToken) String() string {
	return base64.URLEncoding.EncodeToString(t.Selector) + base64.URLEncoding.EncodeToString(t.Validator)
}

func DecodeToken(s string) (UserToken, error) {
	var t UserToken
	selectorEncoded := int(math.Ceil(float64(selectorSize)/3) * 4)
	validatorEncoded := int(math.Ceil(float64(validatorSize)/3) * 4)

	if len(s) != selectorEncoded+validatorEncoded {
		return t, ErrInvalidToken
	}

	var err error
	t.Selector, err = base64.URLEncoding.DecodeString(s[:selectorEncoded])
	if err != nil {
		return t, ErrInvalidToken
	}

	t.Validator, err = base64.URLEncoding.DecodeString(s[selectorEncoded:])
	if err != nil {
		return t, ErrInvalidToken
	}

	return t, nil
}
