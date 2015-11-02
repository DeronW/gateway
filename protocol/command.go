package protocol

import (
	"errors"
	log "github.com/Sirupsen/logrus"
)

func ExpoundPacket(src []byte, ckey *CipherKey) (*Packet, error) {
	bytes, err := decode(src)
	p := &Packet{}

	p.size = int(bytes[0]) + int(bytes[1]<<7) + 1
	if p.size != len(bytes) {
		err = errors.New("packet bytes ivalid, length not matched")
		return p, err
	}

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
			p.Op = bytes2str(reverse(cnt[4:6]))
			p.Params = bytes2str(cnt[6:])
		} else if p.Version == 1 {
			p.Addr = uint32(bytes2int(reverse(cnt[0:4])))
			p.SrcCost = int(bytes2int(reverse(cnt[4:5])))
			p.SrcSeq = int(bytes2int(reverse(cnt[5:6])))
			p.Op = bytes2str(reverse(cnt[8:10]))
			p.Params = bytes2str(cnt[10:])
		}
	} else {
		if p.Version == 0 {
			p.Addr = uint32(bytes2int(reverse(bytes[7:11])))
			p.Op = bytes2str(reverse(bytes[11:12]))
			p.Params = bytes2str(reverse(bytes[17:]))
		} else if p.Version == 1 {
			p.Addr = uint32(bytes2int(reverse(bytes[7:11])))
			p.SrcCost = int(bytes2int(reverse(bytes[11:12])))
			p.SrcSeq = int(bytes2int(reverse(bytes[12:13])))
			p.cmdLength = int(bytes2int(reverse(bytes[13:15])))
			p.Op = bytes2str(reverse(bytes[15:17]))
			// here, should reverse params, but BUG is: it has not do this.
			p.Params = bytes2str(bytes[17 : 17+p.cmdLength-2])
		}
	}

	log.WithFields(log.Fields{
		"packet": p,
	}).Info("find a packet")

	return p, nil
}

func ExpoundCommand(p *Packet) (cmd Command, err error) {
	switch p.Op {
	case "1":
		cmd = &CmdLogin{
			op: p.Op,
		}
	case "2":
		cmd = &CmdLogin{
			op: p.Op,
		}
	}
	return
}
