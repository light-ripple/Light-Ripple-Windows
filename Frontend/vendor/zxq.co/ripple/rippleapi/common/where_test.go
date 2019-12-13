package common

import (
	"reflect"
	"testing"
)

func Test_generateQuestionMarks(t *testing.T) {
	type args struct {
		x int
	}
	tests := []struct {
		name   string
		args   args
		wantQm string
	}{
		{"-1", args{-1}, ""},
		{"0", args{0}, ""},
		{"1", args{1}, "?"},
		{"2", args{2}, "?, ?"},
	}
	for _, tt := range tests {
		if gotQm := generateQuestionMarks(tt.args.x); gotQm != tt.wantQm {
			t.Errorf("%q. generateQuestionMarks() = %v, want %v", tt.name, gotQm, tt.wantQm)
		}
	}
}

func TestWhereClause_In(t *testing.T) {
	type args struct {
		initial string
		fields  []string
	}
	tests := []struct {
		name   string
		fields *WhereClause
		args   args
		want   *WhereClause
	}{
		{
			"simple",
			&WhereClause{},
			args{"users.id", []string{"1", "2", "3"}},
			&WhereClause{"WHERE users.id IN (?, ?, ?)", []interface{}{"1", "2", "3"}, false},
		},
		{
			"withExisting",
			Where("users.username = ?", "Howl").Where("users.xd > ?", "6"),
			args{"users.id", []string{"1"}},
			&WhereClause{
				"WHERE users.username = ? AND users.xd > ? AND users.id IN (?)",
				[]interface{}{"Howl", "6", "1"},
				false,
			},
		},
	}
	for _, tt := range tests {
		w := tt.fields
		if got := w.In(tt.args.initial, tt.args.fields...); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. WhereClause.In() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestWhere(t *testing.T) {
	type args struct {
		clause        string
		passedParam   string
		allowedValues []string
	}
	tests := []struct {
		name string
		args args
		want *WhereClause
	}{
		{
			"simple",
			args{"users.id = ?", "5", nil},
			&WhereClause{"WHERE users.id = ?", []interface{}{"5"}, false},
		},
		{
			"allowed",
			args{"users.id = ?", "5", []string{"1", "3", "5"}},
			&WhereClause{"WHERE users.id = ?", []interface{}{"5"}, false},
		},
		{
			"notAllowed",
			args{"users.id = ?", "5", []string{"0"}},
			&WhereClause{},
		},
	}
	for _, tt := range tests {
		if got := Where(tt.args.clause, tt.args.passedParam, tt.args.allowedValues...); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Where() = %#v, want %#v", tt.name, got, tt.want)
		}
	}
}
