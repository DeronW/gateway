package protocol

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
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

func bytes2int(a []byte) (s int) {
	// 最多转换8个byte
	c := len(a)
	for i := 0; i < c && i < 8; i++ {
		s += int(a[i]) << (uint(c-i-1) * 8)
	}
	return
}

func bytes2str(a []byte) (s string) {
	return fmt.Sprintf("%X", a)
}

func int2byte(i uint64, size int) []byte {
	b := bytes.NewBuffer([]byte{})

	if size == 1 {
		binary.Write(b, binary.LittleEndian, uint8(i))
	} else if size == 2 {
		binary.Write(b, binary.LittleEndian, uint16(i))
	} else if size == 4 {
		binary.Write(b, binary.LittleEndian, uint32(i))
	} else if size == 8 {
		binary.Write(b, binary.LittleEndian, uint64(i))
	} else {
		panic(fmt.Sprintf("fail to convert int(%d), size(%d)", i, size))
	}
	return b.Bytes()
}

func int2str(i uint64, size int) string {
	var s string
	c := int2byte(i, size)
	for m := range c {
		s += string(c[m])
	}
	return s
}
func parseOp(a []byte) string {
	n := bytes2int(reverse(a))
	if n == 1 {
		return "1"
	} else if n == 3 {
		return "3"
	}
	return string(a[0]) + string(a[1])
}

// padding bytes to multiple of aes.BlockSize(16)
func padding16(src []byte, b byte) []byte {
	return append(src, bytes.Repeat([]byte{b}, (16-(len(src)%16))%16)...)
}
