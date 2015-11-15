package protocol

import (
	"fmt"
	"net/url"
)

type Packet struct {
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

type PacketToTeleport struct {
	DeviceAddr        uint32
	Encrypted         bool
	WirelessEncrypted bool
	Op                uint8
	Params            string
}

func (p *Packet) ToRailsURLValues() url.Values {
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

type CipherKey struct {
	UserKeyIndex int
	IV           string
	IvChr        string
	Iv96str      []byte
	EncryptCtr   uint32
	DecryptCtr   uint32
	UserKey      []byte
}

type Command interface {
	GetOp() string
}

type CmdLogin struct {
	op string
}

func (c *CmdLogin) GetOp() string {
	return c.op
}

type CmdCommon struct {
	op string
}

func (c *CmdCommon) GetOp() string {
	return c.op
}
