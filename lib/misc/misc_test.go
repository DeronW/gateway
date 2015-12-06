package misc

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Bytes2int(t *testing.T) {
	Convey("", t, func() {
		x := Bytes2int([]byte{1, 2})
		So(x, ShouldEqual, 258)
	})

	Convey("", t, func() {
		x := Bytes2int([]byte{1, 2, 3, 4, 5, 6, 7, 8})
		So(x, ShouldEqual, 72623859790382856)
	})

	Convey("", t, func() {
		x := Bytes2int([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13})
		So(x, ShouldEqual, 72623859790382856)
	})
}

func Test_Rand8byte(t *testing.T) {
	Convey("", t, func() {
		a := Rand8byte()
		b := Rand8byte()
		c := Rand8byte()

		So(len(a), ShouldEqual, 8)
		So(len(b), ShouldEqual, 8)
		So(len(c), ShouldEqual, 8)

		So(a, ShouldNotEqual, b)
		So(a, ShouldNotEqual, c)
		So(b, ShouldNotEqual, c)
	})
}

func Test_Rc4xor(t *testing.T) {
	Convey("", t, func() {
		a := Rc4xor([]byte{1}, []byte{1})
		So(fmt.Sprintf("%X", a), ShouldEqual, "07")

		b := Rc4xor([]byte{0, 1, 2, 3, 4, 5, 6, 7}, []byte{1})
		So(fmt.Sprintf("%X", b), ShouldEqual, "06090C0D1C252F2E")
	})
}

func Test_BytesXor(t *testing.T) {
	Convey("", t, func() {
		a := BytesXor([]byte{7}, []byte{0})
		So(fmt.Sprintf("%X", a), ShouldEqual, "07")

		b := BytesXor([]byte{7}, []byte{7})
		So(fmt.Sprintf("%X", b), ShouldEqual, "00")
	})
}

func Test_Str2byte(t *testing.T) {
	Convey("", t, func() {
		a, _ := Str2byte("ab")
		So(a[0], ShouldEqual, 171)

		b, _ := Str2byte("abcdef")
		So(b[0], ShouldEqual, 171)
		So(b[1], ShouldEqual, 205)
		So(b[2], ShouldEqual, 239)
	})
}
