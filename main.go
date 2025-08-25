package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"pusher/config"
	"pusher/internal/source"
	"pusher/internal/ws"
	"pusher/pkg/logger"
	"pusher/pkg/redis"
	"pusher/pkg/utils"
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

	redisClient, err := redis.InitClient(&conf.Redis)
	if err != nil {
		logger.GetLogger().Fatalf("failed to connect to redis: %v", err)
	}

	subManager := ws.NewSubscriptionManager(
		source.NewTopicPuller(conf.Source.Redis, source.NewRedisSource(redisClient)),
		source.NewTopicPuller(conf.Source.Kafka, source.NewKafkaSource(&conf.Kafka)),
	)

	server := ws.NewWebsocketServer(&conf.APP, &ws.DefaultHandler{
		SubscriptionManager: subManager,
	})
	if err = server.Run(); err != nil {
		logger.GetLogger().Fatalf("websocket server failed: %v", err)
	}

	utils.WaitForShutdown()

	logger.GetLogger().Info("server shutdown success...")
}
