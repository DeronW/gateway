package protocol

import (
	"fmt"
	"net/url"
)

type PacketReceive struct {
	Encrypted         bool
	WirelessEncrypted bool
	Addr              uint32
	Op                string
	Params            string
	UserKeyIndex      uint16
	SrcCost           int
	SrcSeq            int
	Version           int
	size              int
	cmdLength         int
}

func (p *PacketReceive) ToRailsURLValues() url.Values {
	v := url.Values{}
	v.Set("command[0][encrypted]", fmt.Sprintf("%t", p.Encrypted))
	v.Set("command[0][w_encrypted]", fmt.Sprintf("%t", p.WirelessEncrypted))
	v.Set("command[0][device_addr]", fmt.Sprintf("%d", p.Addr))
	v.Set("command[0][op]", fmt.Sprintf("%s", p.Op))
	v.Set("command[0][params]", fmt.Sprintf("%s", p.Params))
	v.Set("command[0][user_key_index]", fmt.Sprintf("%d", p.UserKeyIndex))
	v.Set("command[0][src_cost]", fmt.Sprintf("%d", p.SrcCost))
	v.Set("command[0][src_seq]", fmt.Sprintf("%d", p.SrcSeq))
	v.Set("command[0][version]", fmt.Sprintf("%d", p.Version))
	v.Set("remote_ip", fmt.Sprintf("%d", p.Version))
	return v
}

func ExpoundPacket(src []byte, ckey *CipherKey) (*PacketReceive, error) {
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
			p.Params = bytes2str(reverse(bytes[14:]))
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
