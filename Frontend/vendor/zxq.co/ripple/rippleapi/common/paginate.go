package common

import "fmt"

// Paginate creates an additional SQL LIMIT clause for paginating.
func Paginate(page, limit string, maxLimit int) string {
	var (
		p = Int(page)
		l = Int(limit)
	)
	if p < 1 {
		p = 1
	}
	if l < 1 {
		l = 50
	}
	if l > maxLimit {
		l = maxLimit
	}
	start := uint(p-1) * uint(l)
	return fmt.Sprintf(" LIMIT %d,%d ", start, l)
}
