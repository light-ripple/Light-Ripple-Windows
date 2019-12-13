package gopher_utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestHexStr2int(t *testing.T) {
	Convey("Convert hex format string to decimal", t, func() {
		hexDecs := map[string]int{
			"1":   1,
			"002": 2,
			"011": 17,
			"0a1": 161,
			"35e": 862,
		}

		for hex, dec := range hexDecs {
			val, err := HexStr2int(hex)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, dec)
		}
	})
}

func TestInt2HexStr(t *testing.T) {
	Convey("Convert decimal to hex format string", t, func() {
		decHexs := map[int]string{
			1:   "1",
			2:   "2",
			17:  "11",
			161: "a1",
			862: "35e",
		}

		for dec, hex := range decHexs {
			val := Int2HexStr(dec)
			So(val, ShouldEqual, hex)
		}
	})
}
