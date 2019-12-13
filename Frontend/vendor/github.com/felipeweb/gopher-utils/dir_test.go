package gopher_utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestIsDir(t *testing.T) {
	Convey("Check if given path is a directory", t, func() {
		Convey("Pass a file name", func() {
			So(IsDir("file.go"), ShouldEqual, false)
		})
		Convey("Pass a directory name", func() {
			So(IsDir("testdata"), ShouldEqual, true)
		})
		Convey("Pass a invalid path", func() {
			So(IsDir("foo"), ShouldEqual, false)
		})
	})
}

func TestCopyDir(t *testing.T) {
	Convey("Items of two slices should be same", t, func() {
		_, err := StatDir("testdata", true)
		So(err, ShouldEqual, nil)

		err = CopyDir("testdata", "testdata2")
		So(err, ShouldEqual, nil)

		_, err = StatDir("testdata2", true)
		os.RemoveAll("testdata2")
		So(err, ShouldEqual, nil)
	})
}
