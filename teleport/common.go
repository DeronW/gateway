package teleport

import (
	"gateway/configs"
	"gateway/protocol/command"
	log "github.com/Sirupsen/logrus"
	"github.com/antonholmquist/jason"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

func Post2Rails(packet *command.Packet, raddr net.Addr) {
	if packet.Op == "1" {
		go post2rails(packet.ToRailsURLValues(), func(bytes []byte) {

			data, err := jason.NewObjectFromBytes(bytes)
			if err != nil {
				log.Info(err)
			}

			e, err := data.GetObject("error")
			if e != nil || err != nil {
				log.Info(e)
				log.Info(err)
				return
			}

			log.Info(data)

			cmd, err := data.GetObject("command")

			if err != nil {
				log.Info(err)
			}

			ctrl, err := data.GetObject("control")

			if err != nil {
				log.Info(err)
			}
		})
	} else {

	}
}

func post2rails(v url.Values, fn func(bytes []byte)) {
	resp, err := http.PostForm(configs.RAILS_SEND_COMMAN_URL, v)
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
