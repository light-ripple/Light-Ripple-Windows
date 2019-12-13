package gopher_utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
	"runtime"
)

func TestDate(t *testing.T) {
	Convey("Convert unix date format to string", t, func() {
		if runtime.GOOS == "darwin" {
			So(Date(1, "DD/MM/YYYY"), ShouldEqual, "31/12/1969")
		} else {
			So(Date(1, "DD/MM/YYYY"), ShouldEqual, "01/01/1970")
		}
	})
	Convey("Convert unix date in string format to humman string", t, func() {
		if runtime.GOOS == "darwin" {
			So(DateS("1", "DD/MM/YYYY"), ShouldEqual, "31/12/1969")
		} else {
			So(DateS("1", "DD/MM/YYYY"), ShouldEqual, "01/01/1970")
		}
	})
	Convey("Convert time object to string", t, func() {
		if runtime.GOOS == "darwin" {
			So(DateT(time.Unix(int64(1), 0), "DD/MM/YYYY"), ShouldEqual, "31/12/1969")
		} else {
			So(DateT(time.Unix(int64(1), 0), "DD/MM/YYYY"), ShouldEqual, "01/01/1970")
		}
	})
}
