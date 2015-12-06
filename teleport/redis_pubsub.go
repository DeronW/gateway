package teleport

import (
	"fmt"
	"gateway/protocol"
	"gopkg.in/redis.v3"
	"sync"
)

var once sync.Once
var redis_client *redis.Client
var pub_ch_0 string = "receive"
var sub_ch_0 string = "should_send"

func publish_packet(teleport int, pk *protocol.PacketReceive) {
	msg := fmt.Sprintf("%d,%d,%s,%s", teleport, pk.Addr, pk.Op, pk.Params)
	err := redis_client.Publish(pub_ch_0, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func subscribe_send() {
	pubsub, err := redis_client.Subscribe(sub_ch_0)
	if err != nil {
		panic(err)
	}
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			continue
		}
		fmt.Println(msg.Payload)
	}
}

func init() {
	once.Do(func() {
		redis_client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		go subscribe_send()
	})
}
