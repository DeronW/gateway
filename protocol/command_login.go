package protocol

import (
	"crypto/aes"
	"fmt"
)

func CommandLoginStepOne(packet *PacketReceive, ck *CipherKey) (
	pk *PacketSend,
	iv []byte,
	user_key []byte,
	user_key_index int,
) {

	nonce1, _ := str2byte(packet.Params)
	nonce2 := rand8byte()
	nonce := append(nonce1, nonce2...)
	private_key, _ := str2byte("55294d59b1f1db94f848fd2364ebc979")
	user_key, _ = str2byte("2d78a9947d265b923c1b55623f13bfb9")
	user_key_index = 0

	out := make([]byte, aes.BlockSize)
	block, _ := aes.NewCipher([]byte(private_key))
	block.Encrypt(out, nonce)

	iv = make([]byte, 12)
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

	pk = &PacketSend{
		Encrypted:         false,
		WirelessEncrypted: false,
		Op:                2,
		Params:            fmt.Sprintf("%X", params),
		Version:           packet.Version,
	}

	return
}
