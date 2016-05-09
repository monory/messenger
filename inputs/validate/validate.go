package validate

import "regexp"

var (
	validUsername = regexp.MustCompile(`^[\w\-]{1,64}$`)
)

func Username(username string) bool {
	return validUsername.MatchString(username)
}

func Password(password string) bool {
	return len(password) <= 256 && len(password) > 0
}
