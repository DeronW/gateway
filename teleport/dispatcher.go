package teleport

import (
	"crypto/aes"
	"errors"
	"fmt"
	"gateway/protocol"
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"runtime/debug"
	"strconv"
	"time"
)

func Dispatch(data []byte, uuid string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(string(debug.Stack()))
		}
	}()

	ckey, err := GetCipherKey(uuid)
	if err != nil {
		return
	}
	packet, cmd, err := protocol.Parse(data, ckey)

	if err != nil {
		return
	}

	switch cmd.GetOp() {
	case "1":
		//Post2RailsLoginCmd(packet, uuid, ckey)
		nonce1, _ := str2byte(packet.Params)
		nonce2 := rand8byte()
		nonce := append(nonce1, nonce2...)
		private_key, _ := str2byte("15294d59b1f1db94f848fd2364ebc979")
		user_key, _ := str2byte("1d78a9947d265b923c1b55623f13bfb9")
		user_key_index := 0

		out := make([]byte, aes.BlockSize)
		block, _ := aes.NewCipher([]byte(private_key))
		block.Encrypt(out, nonce)

		iv := make([]byte, 12)
		for i := 0; i < 12; i++ {
			iv[i] = out[i]
		}
		encrypted_user_key := make([]byte, 16)
		block.Encrypt(encrypted_user_key, bytes_xor(user_key, nonce))

		params := make([]byte, 0, 16)
		params = append(params, out...)
		params = append(params, nonce2...)
		params = append(params, encrypted_user_key...)
		params = append(params, []byte{byte(user_key_index), 0}...)

		//GlobalPool.SetIV(uuid, string(iv), iv_chr)
		GlobalPool.SetIV2(uuid, iv)
		GlobalPool.SetUserKey(uuid, user_key)
		GlobalPool.SetUserKeyIndex(uuid, user_key_index)
		//GlobalPool.SetTeleportAddr(uuid, addr)

		send2teleport(uuid, packet.Version, &protocol.PacketToTeleport{
			Encrypted:         false,
			WirelessEncrypted: false,
			Op:                2,
			Params:            fmt.Sprintf("%X", params),
		}, ckey)

	case "3":
		send2teleport(uuid, packet.Version, &protocol.PacketToTeleport{
			Encrypted:         true,
			WirelessEncrypted: true,
			Op:                4,
			Params:            "",
		}, ckey)
	case "qt":
		log.Info("return time")
	default:
		log.Info("no handler for this command")
	}
	return nil
}

func rand8byte() []byte {
	return []byte("\x9F\x01\xDF\xD8\xD9\x02\x03\x04")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	out := make([]byte, 8)
	for i := 0; i < 8; i++ {
		out[i] = byte(r.Int() & 0xff)
	}
	return out
}

func str2byte(s string) ([]byte, error) {
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

func bytes_xor(a []byte, b []byte) (c []byte) {
	for i := 0; i < len(a); i++ {
		c = append(c, a[i]^b[i])
	}
	return c
}

func reverse(a []byte) []byte {
	c := len(a) - 1
	b := make([]byte, c+1)
	for i := 0; i <= c; i++ {
		b[i] = a[c-i]
	}
	return b
}
