package main

import (
	"fmt"
	"gateway/common"
	"gateway/teleport"
	tcp "github.com/delongw/phantom-tcp"
	"github.com/spf13/viper"
	"time"
)

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/gateway/")
	viper.AddConfigPath("$HOME/.gateway/")
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")

	// set default config value
	viper.SetDefault("tcp.host", "0.0.0.0")
	viper.SetDefault("tcp.port", 6000)
	viper.SetDefault("tcp.send_buffer", 60)
	viper.SetDefault("tcp.receive_buffer", 60)
	viper.SetDefault("tcp.separtor", '*')

	viper.SetDefault("keepalive.deadline", 0)
	viper.SetDefault("keepalive.idle", false)
	viper.SetDefault("keepalive.count", 0)
	viper.SetDefault("keepalive.interval", 0)
	viper.SetDefault("keepalive.auto_close_duration", 0)

	viper.SetDefault("rails.post_url", "http://127.0.0.1/tcp/command")

	viper.SetDefault("log.out", "console")

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error  config  file: %s\n", err))
	}

	go teleport.Run(&tcp.ServerConfig{
		Host:       viper.GetString("tcp.host"),
		Port:       uint32(viper.GetInt("tcp.port")),
		Net:        "tcp",
		SendBuf:    uint32(viper.GetInt("tcp.send_buffer")),
		ReceiveBuf: uint32(viper.GetInt("tcp.receive_buffer")),

		Deadline:          time.Second * viper.GetDuration("keepalive.deadline"),
		KeepAlive:         viper.GetBool("keepalive.enable"),
		KeepAliveIdle:     time.Second * viper.GetDuration("keepalive.idle"),
		KeepAliveCount:    viper.GetInt("keepalive.count"),
		KeepAliveInterval: time.Second * viper.GetDuration("keepalive.interval"),

		Separtor: '*',
	})

	common.HoldOn()
}
