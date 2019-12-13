package common

import "strconv"

// Int converts s to an int. If s in an invalid int, it defaults to 0.
func Int(s string) int {
	r, _ := strconv.Atoi(s)
	return r
}
