package gopher_utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func Test_Pow(t *testing.T) {
	Convey("Power int", t, func() {
		for x := 0; x < 10; x++ {
			for y := 0; y < 8; y++ {
				result := PowInt(x, y)
				result_float := math.Pow(float64(x), float64(y))
				So(result, ShouldEqual, int(result_float))
			}
		}
	})
}
