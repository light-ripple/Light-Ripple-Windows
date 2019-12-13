package locale

import (
	"sort"
	"strconv"
	"strings"
)

// ParseHeader parses an Accept-Language header, and sorts the values.
func ParseHeader(header string) []string {
	if header == "" {
		return nil
	}
	parts := strings.Split(header, ",")

	sort.Slice(parts, func(i, j int) bool {
		return getQuality(parts[i]) > getQuality(parts[j])
	})

	for idx, val := range parts {
		parts[idx] = strings.Replace(strings.SplitN(val, ";q=", 2)[0], "-", "_", 1)
	}

	return parts
}

func getQuality(s string) float32 {
	idx := strings.Index(s, ";q=")
	if idx == -1 {
		return 1
	}

	f, err := strconv.ParseFloat(s[idx+3:], 32)
	if err != nil {
		return 1
	}
	return float32(f)
}
