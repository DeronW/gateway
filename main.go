package main

import (
	"gateway/config"
	"gateway/lib/extra_data"
	"gateway/lib/holdon"
	"gateway/teleport"
	"github.com/codegangsta/cli"
	"os"
)

func start_server(c *cli.Context) {
	go teleport.Run()
	holdon.HoldOn() // in development env, make server blocking
}

func stop_server(c *cli.Context)    {}
func restart_server(c *cli.Context) {}
func server_status(c *cli.Context)  {}

func main() {
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
					println("example:")
					println("    ./gateway ltpk ~/path/data.csv")
					return
				}
				extra_data.ImportTeleportData(c.Args().First())
			},
		},
		{
			Name:  "configtest",
			Usage: "validate config",
			Action: func(c *cli.Context) {
				config.Validate()
			},
		},
	}

	app.Run(os.Args)
}
