package common

import "strconv"

// In picks x if y < x, picks z if y > z, or if none of the previous
// conditions is satisfies, it simply picks y.
func In(x, y, z int) int {
	switch {
	case y < x:
		return x
	case y > z:
		return z
	}
	return y
}

// InString takes y as a string, also allows for a default value should y be
// invalid as a number.
func InString(x int, y string, z, def int) int {
	num, err := strconv.Atoi(y)
	if err != nil {
		return def
	}
	return In(x, num, z)
}
