package gopher_utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIsLetter(t *testing.T) {
	if IsLetter('1') {
		t.Errorf("IsLetter:\n Expect => %v\n Got => %v\n", false, true)
	}

	if IsLetter('[') {
		t.Errorf("IsLetter:\n Expect => %v\n Got => %v\n", false, true)
	}

	if !IsLetter('a') {
		t.Errorf("IsLetter:\n Expect => %v\n Got => %v\n", true, false)
	}

	if !IsLetter('Z') {
		t.Errorf("IsLetter:\n Expect => %v\n Got => %v\n", true, false)
	}
}

func TestReverse(t *testing.T) {
	if Reverse("abcdefg") != "gfedcba" {
		t.Errorf("Reverse:\n Except => %s\n Got =>%s\n", "gfedcba", Reverse("abcdefg"))
	}
}

func Test_ToSnakeCase(t *testing.T) {
	cases := map[string]string{
		"HTTPServer":         "http_server",
		"_camelCase":         "_camel_case",
		"NoHTTPS":            "no_https",
		"Wi_thF":             "wi_th_f",
		"_AnotherTES_TCaseP": "_another_tes_t_case_p",
		"ALL":                "all",
		"_HELLO_WORLD_":      "_hello_world_",
		"HELLO_WORLD":        "hello_world",
		"HELLO____WORLD":     "hello____world",
		"TW":                 "tw",
		"_C":                 "_c",

		"  sentence case  ":                                    "__sentence_case__",
		" Mixed-hyphen case _and SENTENCE_case and UPPER-case": "_mixed_hyphen_case__and_sentence_case_and_upper_case",
	}
	Convey("Convert string into snake case", t, func() {
		for old, new := range cases {
			So(ToSnakeCase(old), ShouldEqual, new)
		}
	})
}
