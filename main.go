package main

import (
	"fmt"
	"gateway/db"
	"gateway/lib/error_handle"
	"gateway/lib/extra_data"
	"gateway/lib/holdon"
	"gateway/teleport"
	"github.com/codegangsta/cli"
	tcp "github.com/delongw/phantom-tcp"
	"github.com/spf13/viper"
	"os"
	"time"
)

func init_config() {
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

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error  config  file: %s\n", err))
	}
}

func teleport_config() *tcp.ServerConfig {
	return &tcp.ServerConfig{
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
	}
}

func start_server(c *cli.Context) {
	db.SetSqlite3Path(viper.GetString("sqlite3"))
	error_handle.SetupRaven(viper.GetString("sentry_dsn"))
	go teleport.Run(
		teleport_config(),
		viper.GetDuration("keepalive.duration"),
		viper.GetString("rails.post_url"),
		false,
	)
	holdon.HoldOn() // in development env, make server blocking
}

func stop_server(c *cli.Context)    {}
func restart_server(c *cli.Context) {}
func server_status(c *cli.Context)  {}

func main() {

	init_config() // this method should be TOP level

	app := cli.NewApp()
	app.Name = "Phantom Gateway Server"
	app.Version = "0.1"
	app.Usage = "Entrance of IoT"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "phantom",
			Email: "delong@huantengsmart.com",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "start TCP server",
			Action: start_server,
		},
		{
			Name:   "stop",
			Usage:  "stop TCP server",
			Action: stop_server,
		},
		{
			Name:   "restart",
			Usage:  "restart TCP server",
			Action: restart_server,
		},
		{
			Name:   "status",
			Usage:  "show server status",
			Action: server_status,
		},
		{
			Name:    "load_teleport_private_key",
			Aliases: []string{"ltpk"},
			Usage:   "load private key from csv file",
			Action: func(c *cli.Context) {
				// csv file format should be:
				// addr,private_key
				// ...
				path := c.Args().First()
				if path == "" {
					println("must supply a file path")
					return
				}
				db.SetSqlite3Path(viper.GetString("sqlite3"))
				extra_data.ImportTeleportData(path)
			},
		},
	}

	app.Run(os.Args)
}
