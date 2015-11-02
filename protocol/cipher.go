package protocol

import (
	"crypto/aes"
	"encoding/base64"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func Decrypt(secret []byte, ckey *CipherKey) (cnt []byte, err error) {
	//user_key
	log.Info("==========")
	log.Info(secret)
	log.Info(ckey)

	if ckey.UserKeyIndex == 0 && ckey.IV == "" {
		err = errors.New("UserKeyIndex or IV must has ONE!")
		return
	}
	//over = len(secret) % 16
	//blocks = len(secret) / 16
	for i := 0; i < len(secret)/16; i++ {
		encryptedBlock = secret[i*16:(i+1)*16] ^ decryptIvStr(ckey)
		//plainBlock =
	}

	return make([]byte, 1), nil
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

	fmt.Printf("%X\n", enc)

	base64Enc := base64.StdEncoding.EncodeToString(
		append(int2byte(uint64((len(enc)+4)/3*4), 2), enc...),
	)
	return fmt.Sprintf("%s*", base64Enc), nil
}

func decryptIvStr(ckey *CipherKey) {
	ctr := ckey.DecryptCtr
	if ctr == 0 {
		ckey.DecryptCtr = (1 << 32) - 1
	} else {
		ckey.DecryptCtr = ctr - 1
	}
}
