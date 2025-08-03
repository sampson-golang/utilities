package strutil

import (
	"regexp"
	"strings"
)

// Squish removes all whitespace from the beginning and end of the string,
// and replaces any sequence of whitespace characters with a single space.
func Squish(s string) string {
	// First trim leading and trailing whitespace
	s = strings.TrimSpace(s)
	// Then replace any sequence of whitespace with a single space
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, " ")
}
