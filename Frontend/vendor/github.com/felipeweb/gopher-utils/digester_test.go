package gopher_utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMd5(t *testing.T) {
	Convey("Should encrypt using MD5", t, func() {
		So(Md5("test"), ShouldEqual, "098f6bcd4621d373cade4e832627b4f6")
	})
}

func TestSha256(t *testing.T) {
	Convey("Should encrypt using SHA256", t, func() {
		So(Sha256("test"), ShouldEqual, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")
	})
}

func TestSha512(t *testing.T) {
	Convey("Should encrypt using SHA512", t, func() {
		So(Sha512("test"), ShouldEqual, "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff")
	})
}