package teleport

import (
	"gateway/configs"
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
			if err != nil {
				handleRailsControl(ctrl, uuid)
			}

			cmd, err := data.GetObject("command")
			if err != nil {
				handleRailsCommand(cmd)
			}

			log.Info(cmd)
			log.Info(ctrl)
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

func handleRailsControl(ctrl *jason.Object, uuid string) {
	iv, err := ctrl.GetInt64("SET_IV")
	if err == nil {
		GlobalPool.SetIV(uuid, iv)
	}
}

func handleRailsCommand(cmd interface{}) {}
