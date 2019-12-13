package common

import "strings"

// SortConfiguration is the configuration of Sort.
type SortConfiguration struct {
	Allowed        []string // Allowed parameters
	Default        string
	DefaultSorting string // if empty, DESC
	Table          string
}

// Sort allows the request to modify how the query is sorted.
func Sort(md MethodData, config SortConfiguration) string {
	if config.DefaultSorting == "" {
		config.DefaultSorting = "DESC"
	}
	if config.Table != "" {
		config.Table += "."
	}
	var sortBy string
	for _, s := range md.Ctx.Request.URI().QueryArgs().PeekMulti("sort") {
		sortParts := strings.Split(strings.ToLower(b2s(s)), ",")
		if contains(config.Allowed, sortParts[0]) {
			if sortBy != "" {
				sortBy += ", "
			}
			sortBy += config.Table + sortParts[0] + " "
			if len(sortParts) > 1 && contains([]string{"asc", "desc"}, sortParts[1]) {
				sortBy += sortParts[1]
			} else {
				sortBy += config.DefaultSorting
			}
		}
	}
	if sortBy == "" {
		sortBy = config.Default
	}
	if sortBy == "" {
		return ""
	}
	return "ORDER BY " + sortBy
}

func contains(a []string, s string) bool {
	for _, el := range a {
		if s == el {
			return true
		}
	}
	return false
}
