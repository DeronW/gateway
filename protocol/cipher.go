package protocol

import (
	"bytes"
	"crypto/aes"
	//"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func Decrypt(secret []byte, ckey *CipherKey) (cnt []byte, err error) {

	if ckey.UserKeyIndex == 0 && ckey.IV == "" {
		err = errors.New("UserKeyIndex or IV must has ONE!")
		return
	}

	paddingCount := (16 - len(secret)%16) % 16
	secret = padding16(secret, 0x00)

	plain := make([]byte, len(secret))

	for i := 0; i < len(secret)/16; i++ {
		str := decryptIvStr(ckey)
		for j := 0; j < 16; j++ {
			plain[i*16+j] = secret[16*i+j] ^ str[j]
			//plain = append(plain, secret[16*i+j]^str[j])
		}
	}

	cnt, err = removeHash(plain[:len(plain)-paddingCount])

	log.Info("222222222222222222222222")
	log.Info(cnt)
	log.Info("222222222222222222222222")

	return
}

func Encrypt(p *PacketToTeleport, version int) (string, error) {
	var enc []byte
	var err error

	if p.Encrypted {

	} else {
		enc, err = spliceNotEncryptedCmd(p, version)
		if err != nil {
			return "", err
		}
	}

	base64Enc := base64.StdEncoding.EncodeToString(
		append(int2byte(uint64((len(enc)+4)/3*4), 2), enc...),
	)
	return fmt.Sprintf("%s*", base64Enc), nil
}

func decryptIvStr(ckey *CipherKey) []byte {
	block, _ := aes.NewCipher(ckey.UserKey[:aes.BlockSize])

	out := make([]byte, aes.BlockSize)
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.LittleEndian, ckey.DecryptCtr)

	block.Encrypt(out, append(ckey.Iv96str, b_buf.Bytes()...))

	if ckey.DecryptCtr == 0 {
		ckey.DecryptCtr = 1<<32 - 1
	} else {
		ckey.DecryptCtr -= 1
	}
	//log.Info("====================")
	//log.Info(append(ckey.Iv96str, b_buf.Bytes()...))
	//log.Info("====================")
	//log.Info(out)
	//log.Info("====================")
	return out
}

func removeHash(src []byte) ([]byte, error) {
	if len(src) < 9 {
		return []byte{}, errors.New("length should not less than 9 " + string(src))
	}

	cal_hash := calculate_hash(src[2:])

	log.Info("-----------------")
	//log.Info("-----------------")
	//log.Info("-----------------")
	//log.Info(cal_hash)
	//log.Info(src)

	if cal_hash[0] == src[0] && cal_hash[1] == src[1] {
		return src[2:], nil
	} else {
		return []byte{}, errors.New("hash is wrong " + string(src))
	}
}

func calculate_hash(src []byte) (out []byte) {

	PublicKey1 := "\x08\x9a\x84\xc5\xa6\xd1\x32\x66\xc4\x9a\xf8\x14\x11\x6e\x63\x13"
	block, _ := aes.NewCipher([]byte(PublicKey1))

	src = append([]byte{0x00, 0x00}, src...)
	src = padding16(src, 0x00)

	out = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(src)/16; i++ {
		t := make([]byte, 16)
		for j := 0; j < 16; j++ {
			t = append(t, src[j]^out[j])
		}
		block.Encrypt(out, t)
	}
	return out[0:2]
}

// padding bytes to multiple of aes.BlockSize(16)
func padding16(src []byte, b byte) []byte {
	size := aes.BlockSize
	return append(src, bytes.Repeat([]byte{b}, (size-(len(src)%size))%size)...)
}
