package protocol

import (
	"encoding/base64"
	"fmt"
)

func decode(src []byte) (dst []byte, err error) {
	//return base64.StdEncoding.DecodeString(string(s))
	dst = make([]byte, len(src))
	_, err = base64.StdEncoding.Decode(dst, src)
	return
}

func reverse(a []byte) []byte {
	c := len(a) - 1
	b := make([]byte, c+1)
	for i := 0; i <= c; i++ {
		b[i] = a[c-i]
	}
	return b
}

func bytes2int(a []byte) (s int64) {
	// 最多转换8个byte
	c := len(a)
	for i := 0; i < c && i < 8; i++ {
		t := int(a[i])
		s += int64(t << (uint(c-i-1) * 8))
	}
	return
}

func bytes2str(a []byte) (s string) {
	if bytes2int(a) == 1 {
		return "1"
	}
	for i := range a {
		s += fmt.Sprintf("%X", a[i])
	}
	return s
}
