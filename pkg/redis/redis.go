package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"pusher/config"
)

func Init(conf *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
