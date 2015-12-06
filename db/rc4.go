package db

import (
	"errors"
	"fmt"
	"gateway/lib/misc"
)

var RC4_KEY string = "huanteng"

func encrypt_private_key(private_key string) string {
	rand_key := misc.Rand8byte()
	key_bytes, _ := misc.Str2byte(private_key)
	secret := misc.Rc4xor(key_bytes, misc.BytesXor(rand_key, []byte(RC4_KEY)))
	return fmt.Sprintf("%X", append(rand_key, secret...))
}

func decrypt_private_key(ekey string) ([]byte, error) {
	if len(ekey) != 48 {
		return nil, errors.New("illegal encrypted private key: " + ekey)
	}
	rand_key, _ := misc.Str2byte(ekey[:16])
	secret, _ := misc.Str2byte(ekey[16:])
	return misc.Rc4xor(secret, misc.BytesXor(rand_key, []byte(RC4_KEY))), nil
}
