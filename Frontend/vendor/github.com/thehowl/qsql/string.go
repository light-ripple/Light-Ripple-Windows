package qsql

import (
	"strconv"
)

// String is just a string, but it implements numerous functions to convert it to various types.
type String string

// String is a shorthand for string(s).
func (s String) String() string {
	return string(s)
}

// Int will convert s to an int. It will return 0 if conversion failed, with no error.
func (s String) Int() int {
	r, _ := strconv.ParseInt(string(s), 10, 0)
	return int(r)
}

// Uint will convert s to an uint. It will return 0 if conversion failed, with no error.
func (s String) Uint() uint {
	r, _ := strconv.ParseUint(string(s), 10, 0)
	return uint(r)
}

// Int64 will convert s to an int64. It will return 0 if conversion failed, with no error.
func (s String) Int64() int64 {
	r, _ := strconv.ParseInt(string(s), 10, 64)
	return r
}

// Uint64 will convert s to an uint64. It will return 0 if conversion failed, with no error.
func (s String) Uint64() uint64 {
	r, _ := strconv.ParseUint(string(s), 10, 64)
	return r
}

// Float64 will convert s to a float64. It will return 0 if conversion failed, with no error.
func (s String) Float64() float64 {
	r, _ := strconv.ParseFloat(string(s), 64)
	return r
}

// Float32 will convert s to a float32. It will return 0 if conversion failed, with no error.
func (s String) Float32() float32 {
	r, _ := strconv.ParseFloat(string(s), 32)
	return float32(r)
}

var truthyValues = [...]string{
	"1",
	"t",
	"true",
	"y",
	"yes",
}

// Bool converts s to a bool.
//
// The following values are true:
//
//  - 1
//  - t
//  - true
//  - y
//  - yes
//
// All other values are false.
// Bool is not case sensitive.
func (s String) Bool() bool {
	for _, el := range truthyValues {
		if string(s) == el {
			return true
		}
	}
	return false
}
