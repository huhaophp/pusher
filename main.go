package main

import (
	"log"
	"pusher/internal/source"
	"pusher/internal/ws"
	"pusher/pkg/redis"
)

func main() {
	client, err := redis.Init()
	if err != nil {
		log.Panicf("new redis err: %+v", err)
	}

	redisSource := source.NewRedisSource(client)

	log.Printf("starting websocket server...")

	topics := []string{"topic1", "topic2", "topic3", "topic4", "topic5", "topic6", "topic7", "topic8", "topic9", "topic10", "topic11", "topic12", "topic13", "topic14", "topic15", "topic16", "topic17", "topic18", "topic19", "topic20", "topic21", "topic22", "topic23", "topic24", "topic25"}

	handler := ws.DefaultHandler{
		SubscriptionManager: ws.NewSubscriptionManager(topics, redisSource),
	}

	server := ws.NewWebsocketServer(ws.DefaultConfig(), &handler)
	if err := server.Run(); err != nil {
		panic(err)
	}

	log.Printf("websocket server started")

	quit := make(chan struct{})
	<-quit
}
