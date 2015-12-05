package db

import ()

var RC4_KEY string = "huanteng"

func EncryptPrivateKey(src string, key string) []byte {
	return []byte{}
}

func DecryptPrivateKey(secret string, key string) []byte {

	key += "0000000000000000"
	key_a := []byte(key[:16])
	key_b := []byte("huantenghuanteng")

	k := make([]byte, 16)
	for i := 0; i < 16; i++ {
		k[i] = key_a[i] ^ key_b[i]
	}
	return []byte{}
}
