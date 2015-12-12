package config

import (
	"fmt"
	ct "github.com/daviddengcn/go-colortext"
	"time"
)

func check(text string, pass bool) {
	//fmt.Print("validate legality of ", text, ": ")
	fmt.Print(text, ": ")
	if pass {
		ct.ChangeColor(ct.Green, false, ct.None, false)
	} else {
		ct.ChangeColor(ct.Red, false, ct.None, false)
	}

	fmt.Println(pass)
	ct.ResetColor()
}

type TCP struct {
	Host              string
	Port              int
	SendBuffer        int
	ReceiveBuffer     int
	AutoCloseDuration time.Duration
	Separtor          byte
}

func (t *TCP) validate() {
	fmt.Println("\nvalidate tcp config...")
	check("tcp.host", t.Host != "")
	check("tcp.port", t.Port != 0)
	check("tcp.send_buffer", t.SendBuffer > 0)
	check("tcp.receive_buffer", t.ReceiveBuffer > 0)
}

type KeepAlive struct {
	Deadline time.Duration
	Enable   bool
	Idle     time.Duration
	Count    int
	Interval time.Duration
}

func (t *KeepAlive) validate() {
	fmt.Println("\nvalidate keep_alive config...")
	fmt.Println("all config has default value, no need to check")
}

type Sqlite3 struct {
	Path string
}

func (t *Sqlite3) validate() {
	fmt.Println("\nvalidate sqlite3 config...")
	check("sqlite3.paht", t.Path != "")
}

type Log struct {
	Out string
}

func (t *Log) validate() {
	fmt.Println("\nvalidate log config...")
	check("log.out", t.Out != "")
}

type Redis struct {
	Host     string
	Port     int
	PubSub   bool
	Password string
}

func (t *Redis) validate() {
	fmt.Println("\nvalidate redis config...")
	check("redis.host", t.Host != "")
	check("redis.port", t.Port != 0)
	//check("redis.pubsub", t.PubSub != nil)
	check("redis.password", t.Password != "")
}

type Sentry struct {
	DSN string
}

func (t *Sentry) validate() {
	fmt.Println("\nvalidate sentry config...")
	check("sentry.dsn", t.DSN != "")
}

type Cipher struct {
}

func (t *Cipher) validate() {
	fmt.Println("\nvalidate cipher config...")
}
