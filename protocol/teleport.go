package protocol

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"gateway/protocol/command"
)

func Parse(src []byte) (packet *command.Packet, cmd command.Command, err error) {
	packet, err = command.ExpoundPacket(src)
	if err != nil {
		return
	}
	cmd, err = command.ExpoundCommand(packet)
	return
}

func Decrypt(s []byte) (cnt []byte) {
	ext := len(s) % 16
	blocks := len(s) / 16
	//plain := ""

	for i := 0; i < blocks; i++ {
		//plain += s[i*16:(i+1)*16] ^ iv_str
	}

	if ext != 0 {

	}
	return
}

func Encrypt(p *command.PacketToTeleport, version int) (string, error) {

	var enc string
	var err error

	if p.Encrypted {

	} else {
		enc, err = spliceNotEncryptedCmd(p, version)
		if err != nil {
			return "", err
		}
	}

	fmt.Println(enc)
	fmt.Println(len(enc))
	enc = fmt.Sprintf("%s%s", int2str(uint64((len(enc)+4)/3*4), 2), enc)

	fmt.Println("99999999999999999")
	fmt.Println(enc)

	base64Enc := base64.StdEncoding.EncodeToString([]byte(enc))

	return fmt.Sprintf("%s*%s", base64Enc, enc), nil
}

func spliceNotEncryptedCmd(p *command.PacketToTeleport, version int) (enc string, err error) {

	var e uint8
	if p.WirelessEncrypted {
		e = 1 << 6
	}
	e += uint8(version & 3)

	if version == 0 {
		enc = fmt.Sprintf(
			"%s\x00\x00\x00\x00%s%s%s%s",
			string(e),
			int2str(uint64(p.DeviceAddr), 4),
			int2str(uint64(p.Op), 2),
			int2str(uint64(len(p.Params)), 2),
			p.Params,
		)
	} else if version == 1 {
		enc = fmt.Sprintf(
			"%s\x00\x00\x00\x00%s\x00\x00%s%s%s",
			string(e),
			int2str(uint64(p.DeviceAddr), 4),
			int2str(uint64(len(p.Params)+2), 2),
			int2str(uint64(p.Op), 2),
			p.Params,
		)
	} else {
		err = errors.New("wrong version")
	}

	return
}

func int2str(i uint64, size int) (s string) {
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

	c := b.Bytes()
	for m := range c {
		s += string(c[m])
	}

	return
}
