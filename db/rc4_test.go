package db

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_encrypt_private_key(t *testing.T) {
	Convey("", t, func() {
		s := encrypt_private_key("01234567890123456789012345678901")
		So(len(s), ShouldEqual, 48)
	})
}

func Test_decrypt_private_key(t *testing.T) {
	Convey("", t, func() {
		a := "01234567890123456789012345678901"
		s := encrypt_private_key(a)
		b, err := decrypt_private_key(s)
		So(err, ShouldEqual, nil)
		So(fmt.Sprintf("%X", b), ShouldEqual, a)
	})
}
