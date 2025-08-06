package main

import (
	"flag"
	"fmt"
	"pusher/config"
	"pusher/internal/source"
	"pusher/internal/ws"
	"pusher/pkg/logger"
	"pusher/pkg/redis"
)

var (
	path = flag.String("cfg", "./config/config.yaml", "config file path")
)

func main() {
	flag.Parse()

	conf, err := config.LoadConfig(*path)
	if err != nil {
		panic(fmt.Sprintf("load config failed: %v", err))
	}

	logger.Init(&conf.Logger)

	client, err := redis.Init(&conf.Redis)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}

	manager := ws.NewSubscriptionManager(
		conf.Source.Redis,
		source.NewRedisSource(client),
	)

	handler := ws.DefaultHandler{SubscriptionManager: manager}

	server := ws.NewWebsocketServer(&conf.APP, &handler)
	if err := server.Run(); err != nil {
		logger.Fatalf("failed to start websocket server: %v", err)
	}

	logger.Info("websocket server started and running...")

	select {}
}
