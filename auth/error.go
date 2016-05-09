package auth

type AuthError struct {
	e string
}

func (err AuthError) Error() string {
	return err.e
}

var (
	ErrUserExists       = AuthError{"auth: user exists"}
	ErrUsernameNotFound = AuthError{"auth: username not found"}
	ErrInvalidPassword  = AuthError{"auth: invalid password"}
	ErrInvalidToken     = AuthError{"auth: invalid token"}
)
