package common

import (
	"unicode"
)

// SanitiseString removes all control codes from a string.
func SanitiseString(s string) string {
	n := make([]rune, 0, len(s))
	for _, c := range s {
		if c == '\n' || !unicode.Is(unicode.Other, c) {
			n = append(n, c)
		}
	}
	return string(n)
}
