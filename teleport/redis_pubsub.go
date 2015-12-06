package teleport

import (
	"fmt"
	"gopkg.in/redis.v3"
	"sync"
)

var once sync.Once
var redis_client *redis.Client
var pub_ch_0 string = "receive"
var sub_ch_0 string = "should_send"

func redis_publish(msg string) {
	err := redis_client.Publish(pub_ch_0, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func redis_subscribe() {
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
		go redis_subscribe()
	})
}
