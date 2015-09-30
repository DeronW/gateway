package teleport

import (
	"gateway/protocol"
	"gateway/protocol/command"
	log "github.com/Sirupsen/logrus"
	"github.com/antonholmquist/jason"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Post2Rails(packet *command.Packet, uuid string) {
	defer func() {
		if err := recover(); err != nil {
			log.Info("post to rails got error")
			log.Info(err)
		}
	}()

	if packet.Op == "1" {
		go post2rails(packet.ToRailsURLValues(), func(bytes []byte) {
			data, _ := jason.NewObjectFromBytes(bytes)
			e, _ := data.GetObject("error")
			if e != nil {
				log.Info(e)
				return
			}

			log.Info(data)
			ctrl, err := data.GetObject("control")
			if err == nil {
				handleRailsControl(uuid, ctrl)
			}

			cmd, err := data.GetObject("command")
			if err == nil {
				handleRailsCommand(uuid, packet.Version, cmd)
			}
		})
	} else {

	}
}

func post2rails(v url.Values, fn func(bytes []byte)) {
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
	if err == nil {
		iv, _ := iv.Number()
		GlobalPool.SetIV(uuid, string(iv))
	}

	uk, err := ctrl.GetString("SET_USER_KEY")
	if err == nil {
		GlobalPool.SetUserKey(uuid, uk)
	}

	uki, err := ctrl.GetInt64("SET_USER_KEY_INDEX")
	if err == nil {
		GlobalPool.SetUserKeyIndex(uuid, int(uki))
	}
}

func handleRailsCommand(uuid string, version int, cmd *jason.Object) {
	p := &command.PacketToTeleport{}
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
