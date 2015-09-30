package common

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSpec(t *testing.T) {
	Convey("Test", t, func() {
		x := 1

		Convey("Value", func() {
			x++

			Convey("equal 2", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}
