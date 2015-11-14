package protocol

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	//log "github.com/Sirupsen/logrus"
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
		bs := decryptIvStr(ckey)
		for j := 0; j < 16; j++ {
			plain[i*16+j] = secret[16*i+j] ^ bs[j]
		}
	}
	cnt, err = removeHash(plain[:len(plain)-paddingCount])
	return
}

func Encrypt(p *PacketToTeleport, version int, ckey *CipherKey) (string, error) {
	var enc []byte
	var err error

	if p.Encrypted {
		enc, err = spliceEncryptedCmd(p, version, ckey)
	} else {
		enc, err = spliceNotEncryptedCmd(p, version)
	}

	if err != nil {
		return "", err
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
	return out
}

func encryptIvStr(ckey *CipherKey) []byte {
	out := make([]byte, aes.BlockSize)
	block, _ := aes.NewCipher(ckey.UserKey[:aes.BlockSize])
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, ckey.EncryptCtr)
	block.Encrypt(out, append(ckey.Iv96str, buf.Bytes()...))
	ckey.EncryptCtr = (ckey.EncryptCtr + 1) & (1<<32 - 1)
	return out
}

func removeHash(src []byte) ([]byte, error) {
	if len(src) < 9 {
		return []byte{}, errors.New("length should not less than 9 " + string(src))
	}

	cal_hash := calculate_hash(src[2:])
	if cal_hash[0] == src[0] && cal_hash[1] == src[1] {
		return src[2:], nil
	} else {
		return []byte{}, errors.New("hash is wrong " + bytes2str(src))
	}
}

func calculate_hash(src []byte) []byte {
	PublicKey1 := "\x08\x9a\x84\xc5\xa6\xd1\x32\x66\xc4\x9a\xf8\x14\x11\x6e\x63\x13"
	block, _ := aes.NewCipher([]byte(PublicKey1))

	src = append([]byte{0x00, 0x00}, src...)
	src = padding16(src, 0x00)

	out := make([]byte, 16)
	for i := 0; i < len(src)/16; i++ {
		t := make([]byte, 16)
		for j := 0; j < 16; j++ {
			t[j] = src[i*16+j] ^ out[j]
		}
		block.Encrypt(out, t)
	}
	return out[0:2]
}

func spliceEncryptedCmd(p *PacketToTeleport, version int, ckey *CipherKey) (enc []byte, err error) {

	var src []byte
	if version == 0 {
		src = append(src, int2byte(uint64(p.DeviceAddr), 4)...)
		src = append(src, int2byte(uint64(p.Op), 2)...)
		src = append(src, params_size_v0(p.Params)...)
		params, err := str2byte(p.Params)
		if err != nil {
			return enc, err
		}
		src = append(src, params...)
	} else if version == 1 {
		src = append(src, int2byte(uint64(p.DeviceAddr), 4)...)
		src = append(src, 0x00, 0x00)
		src = append(src, params_size_v1(p.Params)...)
		src = append(src, int2byte(uint64(p.Op), 2)...)
		params, err := str2byte(p.Params)
		if err != nil {
			return enc, err
		}
		src = append(src, params...)
	} else {
		return enc, errors.New("Wrong command version")
	}
	encryption := int(version&3) + 1<<7
	if p.Encrypted {
		encryption += 1 << 6
	}
	enc = append(enc, byte(encryption))
	enc = append(enc, int2byte(uint64(ckey.UserKeyIndex), 2)...)
	secret, err := encrypt_plain_cmd(src, ckey)
	if err != nil {
		return enc, err
	}
	enc = append(enc, secret...)
	return
}

func spliceNotEncryptedCmd(p *PacketToTeleport, version int) (enc []byte, err error) {

	e := int(version & 3)
	if p.WirelessEncrypted {
		e += 1 << 6
	}

	enc = append(enc, byte(e))

	if version == 0 {
		enc = append(enc, 0x00, 0x00, 0x00, 0x00)
		enc = append(enc, int2byte(uint64(p.DeviceAddr), 4)...)
		enc = append(enc, int2byte(uint64(p.Op), 2)...)
		enc = append(enc, params_size_v0(p.Params)...)
		params, err := str2byte(p.Params)
		if err != nil {
			return enc, err
		}
		enc = append(enc, params...)
	} else if version == 1 {
		enc = append(enc, 0x00, 0x00, 0x00, 0x00)
		enc = append(enc, int2byte(uint64(p.DeviceAddr), 4)...)
		enc = append(enc, 0x00, 0x00)
		enc = append(enc, params_size_v1(p.Params)...)
		enc = append(enc, int2byte(uint64(p.Op), 2)...)

		params, err := str2byte(p.Params)

		if err != nil {
			return enc, err
		}
		enc = append(enc, params...)
	} else {
		err = errors.New("wrong version")
		return
	}

	return
}

func params_size_v0(params string) []byte {
	return int2byte(uint64(len(params)), 1)
}

func params_size_v1(params string) []byte {
	return int2byte(uint64(len(params)/2+2), 2)
}

func encrypt_plain_cmd(plain []byte, ckey *CipherKey) ([]byte, error) {
	hash := calculate_hash(plain)
	plain = append(hash, plain...)

	paddingCount := (16 - len(plain)%16) % 16
	plain = padding16(plain, 0x00)

	secret := make([]byte, len(plain))

	for i := 0; i < len(plain)/16; i++ {
		bs := encryptIvStr(ckey)
		for j := 0; j < 16; j++ {
			secret[i*16+j] = plain[16*i+j] ^ bs[j]
		}
	}

	return secret[:len(plain)-paddingCount], nil
}
