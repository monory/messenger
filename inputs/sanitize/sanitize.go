package sanitize

import (
	"strings"
)

func Username(username string) string {
	return strings.ToLower(username)
}
