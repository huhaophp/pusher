package source

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pusher/internal/types"
	"pusher/pkg/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSource struct {
	redis      *redis.Client
	enableMock bool
}

func NewRedisSource(redis *redis.Client) *RedisSource {
	return &RedisSource{
		redis:      redis,
		enableMock: true, // 默认关闭模拟
	}
}

func (r *RedisSource) PullMessage(ctx context.Context, topic string, handler func(data *types.Data)) error {
	if r.enableMock {
		r.startMockProducer(topic)
	}

	pubSub := r.redis.Subscribe(ctx, topic)

	if _, err := pubSub.Receive(ctx); err != nil {
		return fmt.Errorf("subscribe error: %w", err)
	}

	logger.GetLogger().WithField("topic", topic).Info("redis consumer started")

	ch := pubSub.Channel(
		redis.WithChannelSize(500),
		redis.WithChannelHealthCheckInterval(30*time.Second),
	)

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return errors.New("redis channel closed")
			}
			var data types.Data
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				logger.GetLogger().Warnf("unmarshal error: %v", err)
				continue
			}
			data.Meta.ReceiveTime = time.Now()
			handler(&data)
		case <-ctx.Done():
			logger.GetLogger().Info("context cancelled")
			return nil
		}
	}
}

func (r *RedisSource) startMockProducer(topic string) {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 50)
		defer ticker.Stop()
		for range ticker.C {
			payload := fmt.Sprintf("mock data at %s", time.Now().Format(time.RFC3339))
			typeList := []string{"test", "test2"}
			for _, typ := range typeList {
				data, _ := json.Marshal(types.Data{
					Topic:   topic,
					Type:    typ,
					Payload: payload,
				})
				if err := r.redis.Publish(context.Background(), topic, data).Err(); err != nil {
					logger.GetLogger().Warnf("mock publish failed: %v", err)
				}
			}
		}
	}()

}
