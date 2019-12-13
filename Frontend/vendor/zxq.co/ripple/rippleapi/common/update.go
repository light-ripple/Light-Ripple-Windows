package common

import (
	"reflect"
	"strings"
)

// UpdateQuery is simply an SQL update query,
// that can be built upon passed parameters.
type UpdateQuery struct {
	fields     []string
	Parameters []interface{}
}

// Add adds a new field with correspective value to UpdateQuery
func (u *UpdateQuery) Add(field string, value interface{}) *UpdateQuery {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return u
	}
	if s, ok := value.(string); ok && s == "" {
		return u
	}
	u.fields = append(u.fields, field+" = ?")
	u.Parameters = append(u.Parameters, value)
	return u
}

// Fields retrieves the fields joined by a comma.
func (u *UpdateQuery) Fields() string {
	return strings.Join(u.fields, ", ")
}
