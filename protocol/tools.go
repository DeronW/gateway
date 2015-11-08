package protocol

import (
	"encoding/base64"
	"fmt"
)

func decode(src []byte) ([]byte, error) {
	dst := make([]byte, len(src))
	size, err := base64.StdEncoding.Decode(dst, src)
	return dst[:size], err
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
	return fmt.Sprintf("%X", a)
}

func parseOp(a []byte) string {
	if bytes2int(a) == 1 {
		return "1"
	}
	return string(a[0]) + string(a[1])
}
