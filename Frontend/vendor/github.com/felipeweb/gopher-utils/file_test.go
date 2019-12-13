package gopher_utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIsFile(t *testing.T) {
	if !IsFile("file.go") {
		t.Errorf("IsExist:\n Expect => %v\n Got => %v\n", true, false)
	}

	if IsFile("testdata") {
		t.Errorf("IsExist:\n Expect => %v\n Got => %v\n", false, true)
	}

	if IsFile("files.go") {
		t.Errorf("IsExist:\n Expect => %v\n Got => %v\n", false, true)
	}
}

func TestIsExist(t *testing.T) {
	Convey("Check if file or directory exists", t, func() {
		Convey("Pass a file name that exists", func() {
			So(IsExist("file.go"), ShouldEqual, true)
		})
		Convey("Pass a directory name that exists", func() {
			So(IsExist("testdata"), ShouldEqual, true)
		})
		Convey("Pass a directory name that does not exist", func() {
			So(IsExist(".hg"), ShouldEqual, false)
		})
	})
}
