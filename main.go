package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"pusher/config"
	"pusher/internal/source"
	"pusher/internal/ws"
	"pusher/pkg/logger"
	"pusher/pkg/redis"
	"pusher/pkg/utils"

	"golang.org/x/sync/errgroup"
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

	subscriptionManager := ws.NewSubscriptionManager(
		source.NewTopicPuller(conf.Source.Redis, source.NewRedisSource(redisClient)),
		source.NewTopicPuller(conf.Source.Kafka, source.NewKafkaSource(&conf.Kafka)),
	)

	group := errgroup.Group{}
	group.Go(func() error {
		if !conf.PProf.Enable {
			return nil
		}
		if err = http.ListenAndServe(fmt.Sprintf(":%s", conf.PProf.Port), nil); err != nil {
			return fmt.Errorf("failed to start pprof server: %v", err)
		}
		return nil
	})
	group.Go(func() error {
		server := ws.NewWebsocketServer(&conf.APP, &ws.DefaultHandler{
			SubscriptionManager: subscriptionManager,
		})
		if err = server.Run(); err != nil {
			return fmt.Errorf("websocket server failed: %v", err)
		}
		return nil
	})

	err = group.Wait()
	if err != nil {
		logger.GetLogger().Fatalf("error occurred: %v", err)
	}

	utils.WaitForShutdown()

	logger.GetLogger().Info("server shutdown success...")
}
