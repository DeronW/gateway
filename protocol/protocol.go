package protocol

import (
	"crypto/aes"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func Parse(src []byte, key *CipherKey) (pk *PacketReceive, cmd Command, err error) {
	pk, err = expound_packet(src, key)
	if err != nil {
		return
	}
	cmd, err = expound_command(pk)

	log.WithFields(log.Fields{
		"packet":  pk,
		"command": cmd,
	}).Info("receive a packet")

	return
}

func LoginStepOne(packet *PacketReceive, ck *CipherKey) (
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

func LoginStepThree() {}

func expound_packet(src []byte, ckey *CipherKey) (*PacketReceive, error) {
	p := &PacketReceive{}
	bytes, err := decode(src)
	if err != nil {
		return p, err
	}

	p.size = int(bytes[0]) + int(bytes[1]<<7) + 1
	encr := bytes[2]

	p.Encrypted = (encr&128 != 0)
	p.WirelessEncrypted = encr&64 == 1
	p.Version = int(encr & 3)

	if p.Encrypted {
		uki := reverse(bytes[3:5])
		ckey.UserKeyIndex = int(bytes2int(reverse(uki)))
		cnt, err := Decrypt(bytes[5:], ckey)
		if err != nil {
			return nil, err
		}
		if p.Version == 0 {
			p.Addr = uint32(bytes2int(reverse(cnt[0:4])))
			p.Op = parseOp(cnt[4:6])
			p.Params = bytes2str(cnt[6:])
		} else if p.Version == 1 {
			p.Addr = uint32(bytes2int(reverse(cnt[0:4])))
			p.SrcCost = int(bytes2int(reverse(cnt[4:5])))
			p.SrcSeq = int(bytes2int(reverse(cnt[5:6])))
			p.Op = parseOp(cnt[8:10])
			p.Params = bytes2str(cnt[10:])
		}
	} else {
		if p.Version == 0 {
			p.Addr = uint32(bytes2int(reverse(bytes[7:11])))
			p.Op = parseOp(bytes[11:12])
			p.Params = bytes2str(reverse(bytes[17:]))
		} else if p.Version == 1 {
			p.Addr = uint32(bytes2int(reverse(bytes[7:11])))
			p.SrcCost = int(bytes2int(reverse(bytes[11:12])))
			p.SrcSeq = int(bytes2int(reverse(bytes[12:13])))
			p.cmdLength = int(bytes2int(reverse(bytes[13:15])))
			p.Op = parseOp(bytes[15:17])
			p.Params = bytes2str(bytes[17 : 17+p.cmdLength-2])
		}
	}
	return p, nil
}

func expound_command(p *PacketReceive) (cmd Command, err error) {
	switch p.Op {
	case "1":
		cmd = &CmdLogin{
			op: p.Op,
		}
	case "3":
		cmd = &CmdLogin{
			op: p.Op,
		}
	default:
		cmd = &CmdCommon{
			op: p.Op,
		}
	}
	return
}
