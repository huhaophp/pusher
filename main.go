package main

import (
	"flag"
	"fmt"
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
		logger.Fatalf("failed to connect to redis: %v", err)
	}

	redisTopicPuller := source.NewTopicPuller(conf.Source.Redis, source.NewRedisSource(redisClient))
	kafkaTopicPuller := source.NewTopicPuller(conf.Source.Kafka, source.NewKafkaSource(&conf.Kafka))

	subscriptionManager := ws.NewSubscriptionManager(redisTopicPuller, kafkaTopicPuller)

	server := ws.NewWebsocketServer(&conf.APP, &ws.DefaultHandler{
		SubscriptionManager: subscriptionManager,
	})
	if err := server.Run(); err != nil {
		logger.Fatalf("failed to start websocket server: %v", err)
	}

	logger.Info("websocket server started and running...")

	utils.WaitForShutdown()

	logger.Infof("websocket server shutdown success...")
}
