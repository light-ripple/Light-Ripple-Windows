package common

import (
	"strings"
)

// SafeUsername makes a string lowercase and replaces all spaces with
// underscores.
func SafeUsername(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "_", -1)
}
