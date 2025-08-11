package source

import (
	"context"
	"encoding/json"
	"fmt"
	"pusher/internal/types"
	"pusher/pkg/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSource struct {
	redis      *redis.Client // redis 客户端
	enableMock bool          // 是否开启模拟数据
}

func NewRedisSource(redis *redis.Client) *RedisSource {
	return &RedisSource{
		redis:      redis,
		enableMock: true,
	}
}

// PullMessage 拉取消息
func (r *RedisSource) PullMessage(ctx context.Context, topic string, handler func(data *types.Data)) error {
	if r.enableMock {
		go r.mockMessage(topic)
	}
	subscribe := r.redis.Subscribe(ctx, topic)
	defer subscribe.Close()
	if _, err := subscribe.Receive(ctx); err != nil {
		return fmt.Errorf("receive error: %w", err)
	}
	logger.GetLogger().WithField("indexer", "RedisSource").Infof("redis consumer started for topic: %s", topic)
	for {
		select {
		case msg := <-subscribe.Channel():
			if msg == nil {
				logger.GetLogger().Info("msg is nil")
				continue
			}
			var data types.Data
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				logger.GetLogger().Infof("json umarshal failed, err: %+v", err)
				continue
			}
			data.Meta.ReceiveTime = time.Now()
			handler(&data)
		case <-ctx.Done():
			logger.GetLogger().Info("context done")
			return nil
		}
	}
}

// mockMessage 模拟消息
func (r *RedisSource) mockMessage(topic string) {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		var customPayload struct {
			Msg string `json:"msg"`
		}
		customPayload.Msg = fmt.Sprintf("hello world, time: %+v", time.Now().Format("2006-01-02 15:04:05"))
		payload, _ := json.Marshal(customPayload)
		data, _ := json.Marshal(map[string]any{
			"topic":   topic,
			"type":    "test",
			"payload": string(payload),
		})
		r.redis.Publish(context.Background(), topic, data)
		time.Sleep(time.Second)
	}
}
