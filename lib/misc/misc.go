package misc

import (
	"crypto/rc4"
	"errors"
	"math/rand"
	"strconv"
	"time"
)

// 最多转换8个byte
func Bytes2int(a []byte) (s int) {
	c := len(a)
	for i := 0; i < c && i < 8; i++ {
		s += int(a[i]) << (uint(c-i-1) * 8)
	}
	return
}

func Rand8byte() []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	out := make([]byte, 8)
	for i := 0; i < 8; i++ {
		out[i] = byte(r.Int() & 0xff)
	}
	return out
}

func Rc4xor(src []byte, k []byte) []byte {
	s := src
	c, err := rc4.NewCipher(k)
	if err != nil {
		panic(err)
	}
	out := make([]byte, len(s))
	c.XORKeyStream(out, s)
	return out
}

func BytesXor(a []byte, b []byte) (c []byte) {
	len_a := len(a)
	len_b := len(b)
	for i := 0; i < len_a && i < len_b; i++ {
		c = append(c, a[i]^b[i])
	}
	return c
}

func Str2byte(s string) ([]byte, error) {
	var b []byte
	length := len(s)
	if length%2 != 0 {
		return make([]byte, 0), errors.New("params is not odd")
	}
	for i := 0; i < length; i += 2 {
		a, err := strconv.ParseUint(string(s[i:i+2]), 16, 8)
		if err != nil {
			return b, err
		}
		b = append(b, byte(a))
	}
	return b, nil
}
