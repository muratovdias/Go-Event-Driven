package main

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
	"os"
)

func main() {
	logger := watermill.NewStdLogger(false, false)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}
	defer func(publisher *redisstream.Publisher) {
		err := publisher.Close()
		if err != nil {
			fmt.Printf("unavailable close publisher, %s", err.Error())
			os.Exit(1)
		}
	}(publisher)

	msg := message.NewMessage(watermill.NewUUID(), []byte("50"))
	err = publisher.Publish("progress", msg)
	if err != nil {
		panic(err)
	}
	msg = message.NewMessage(watermill.NewUUID(), []byte("100"))
	err = publisher.Publish("progress", msg)
	if err != nil {
		panic(err)
	}

}
