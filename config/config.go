package config

import (
	"github.com/spf13/viper"
	"sync"
	"time"
)

var once sync.Once

func Validate() {
	GetTCPConfig().validate()
	GetKeepAliveCfg().validate()
	GetSqlite3Cfg().validate()
	GetLogCfg().validate()
	GetRedisCfg().validate()
	GetSentryCfg().validate()
	GetCipherCfg().validate()
}

func GetTCPConfig() *TCP {
	return &TCP{
		Host:              viper.GetString("tcp.host"),
		Port:              viper.GetInt("tcp.port"),
		SendBuffer:        viper.GetInt("tcp.send_buffer"),
		ReceiveBuffer:     viper.GetInt("tcp.receive_buffer"),
		AutoCloseDuration: time.Second * viper.GetDuration("tcp.auto_close_duration"),
		Separtor:          '*',
	}
}

func GetSqlite3Cfg() *Sqlite3 {
	return &Sqlite3{viper.GetString("sqlite3")}
}

func GetKeepAliveCfg() *KeepAlive {
	return &KeepAlive{
		Deadline: time.Second * viper.GetDuration("keepalive.deadline"),
		Enable:   viper.GetBool("keepalive.enable"),
		Idle:     time.Second * viper.GetDuration("keepalive.idle"),
		Count:    viper.GetInt("keepalive.count"),
		Interval: time.Second * viper.GetDuration("keepalive.interval"),
	}
}

func GetLogCfg() *Log {
	return &Log{
		Out: viper.GetString("log.out"),
	}
}

func GetRedisCfg() *Redis {
	return &Redis{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetInt("redis.port"),
		PubSub:   viper.GetBool("redis.pubsub"),
		Password: viper.GetString("redis.password"),
	}
}

func GetSentryCfg() *Sentry {
	return &Sentry{
		DSN: viper.GetString("sentry.dsn"),
	}
}

func GetCipherCfg() *Cipher {
	return &Cipher{}
}

func set_default_cfg() {
	// set default config value
	viper.SetDefault("tcp.host", "0.0.0.0")
	viper.SetDefault("tcp.port", 6000)
	viper.SetDefault("tcp.send_buffer", 60)
	viper.SetDefault("tcp.receive_buffer", 60)

	viper.SetDefault("keepalive.deadline", 0)
	viper.SetDefault("keepalive.idle", false)
	viper.SetDefault("keepalive.count", 0)
	viper.SetDefault("keepalive.interval", 0)
	viper.SetDefault("keepalive.auto_close_duration", 0)

	viper.SetDefault("rails.post_url", "http://127.0.0.1/tcp/command")
	viper.SetDefault("log.out", "console")
}

func init() {
	once.Do(func() {
		viper.SetConfigName("config")
		//viper.AddConfigPath("/etc/gateway/")
		//viper.AddConfigPath("$HOME/.gateway/")
		viper.AddConfigPath("./")
		viper.SetConfigType("yaml")

		if err := viper.ReadInConfig(); err != nil {
			println("Fatal error config file: %s\n", err)
			return
		}

		set_default_cfg()
	})
}
