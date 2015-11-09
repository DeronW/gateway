package teleport

import (
	"gateway/protocol"
	log "github.com/Sirupsen/logrus"
	"github.com/antonholmquist/jason"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func Post2RailsLoginCmd(packet *protocol.Packet, uuid string) {
	defer func() {
		if err := recover(); err != nil {
			log.Info("post to rails got error")
			log.Info(err)
		}
	}()

	if packet.Op == "1" {
		go post2rails(packet.ToRailsURLValues(), func(body []byte) {
			data, _ := jason.NewObjectFromBytes(body)
			e, _ := data.GetObject("error")
			if e != nil {
				log.Info(e)
				return
			}

			ctrl, err := data.GetObject("control")
			if err == nil {
				handleRailsControl(uuid, ctrl)
			}

			cmd, err := data.GetObject("command")
			if err == nil {
				handleRailsCommand(uuid, packet.Version, cmd)
			}
		})
	} else if packet.Op == "3" {
		go post2rails(packet.ToRailsURLValues(), func(body []byte) {
			data, _ := jason.NewObjectFromBytes(body)
			e, _ := data.GetObject("error")
			if e != nil {
				log.Info(e)
				return
			}

			ctrl, err := data.GetObject("control")
			if err == nil {
				handleRailsControl(uuid, ctrl)
			}
		})
	}
}

func post2rails(v url.Values, fn func(body []byte)) {
	resp, err := http.PostForm(RailsPostUrl, v)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		resp.Body.Close()
		if err := recover(); err != nil {
			log.Info(err)
			log.Info(string(body))
		}
	}()

	if err != nil {
		log.Fatal(err)
	}

	body = body[1 : len(body)-1]
	fn([]byte(body))
}

func handleRailsControl(uuid string, ctrl *jason.Object) {
	iv, err := ctrl.GetValue("SET_IV")
	iv_chr, err := ctrl.GetString("SET_IV_CHR")
	if err == nil {
		iv, _ := iv.Number()
		GlobalPool.SetIV(uuid, string(iv), iv_chr)
	}

	uk, err := ctrl.GetString("SET_USER_KEY")
	if err == nil {
		buk := make([]byte, 0, 16)
		for i := 0; i < 32; i += 2 {
			b, _ := strconv.ParseInt(uk[i:i+2], 16, 16)
			buk = append(buk, byte(b))
		}
		GlobalPool.SetUserKey(uuid, buk)
	}

	uki, err := ctrl.GetInt64("SET_USER_KEY_INDEX")
	if err == nil {
		GlobalPool.SetUserKeyIndex(uuid, int(uki))
	}

	addr, err := ctrl.GetInt64("SET_TELEPORT_ADDR")
	if err == nil {
		GlobalPool.SetTeleportAddr(uuid, addr)
	}
}

func handleRailsCommand(uuid string, version int, cmd *jason.Object) {
	p := &protocol.PacketToTeleport{}
	addr, _ := cmd.GetInt64("device_addr")
	p.DeviceAddr = uint16(addr)

	p.Encrypted, _ = cmd.GetBoolean("encrypted")

	op, _ := cmd.GetInt64("op")
	p.Op = uint8(op)
	p.Params, _ = cmd.GetString("params")
	p.WirelessEncrypted, _ = cmd.GetBoolean("w_encrypted")

	enc, err := protocol.Encrypt(p, version)

	if err != nil {
		log.Error(err)
		return
	}

	GlobalPool.Send(uuid, enc)
}
