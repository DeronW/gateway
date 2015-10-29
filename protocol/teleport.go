package protocol

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"gateway/protocol/command"
	"strconv"
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

func spliceNotEncryptedCmd(p *command.PacketToTeleport, version int) ([]byte, error) {

	var enc []byte
	e := int(version & 3)
	if p.WirelessEncrypted {
		e += 1 << 6
	}

	enc = append(enc, byte(e))
	paramsSize := len(p.Params)/2 + 2

	if version == 0 {
		enc = append(enc, 0x00, 0x00, 0x00, 0x00)
		enc = append(enc, int2byte(uint64(p.DeviceAddr), 4)...)
		enc = append(enc, int2byte(uint64(p.Op), 2)...)
		enc = append(enc, int2byte(uint64(paramsSize), 2)...)
		params, err := str2byte(p.Params)
		if err != nil {
			return make([]byte, 0), err
		}
		enc = append(enc, params...)
	} else if version == 1 {
		enc = append(enc, 0x00, 0x00, 0x00, 0x00)
		enc = append(enc, int2byte(uint64(p.DeviceAddr), 4)...)
		enc = append(enc, 0x00, 0x00)
		enc = append(enc, int2byte(uint64(paramsSize), 2)...)
		enc = append(enc, int2byte(uint64(p.Op), 2)...)

		params, err := str2byte(p.Params)

		if err != nil {
			return make([]byte, 0), err
		}
		enc = append(enc, params...)
	} else {
		return []byte{}, errors.New("wrong version")
	}

	return enc, nil
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
