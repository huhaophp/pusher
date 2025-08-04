package source

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"pusher/internal/model"
	"time"
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

// PullMessage 获取消息.
func (r *RedisSource) PullMessage(ctx context.Context, topic string, handler func(data *model.Data)) error {
	if r.enableMock {
		go r.mockMessage(topic)
	}
	subscribe := r.redis.Subscribe(ctx, topic)
	defer subscribe.Close()
	if _, err := subscribe.Receive(ctx); err != nil {
		return fmt.Errorf("receive error: %w", err)
	}
	log.Printf("RedisSource subscribe %s success", topic)
	for {
		select {
		case msg := <-subscribe.Channel():
			if msg == nil {
				log.Printf("msg is nil")
				continue
			}
			var data model.Data
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				log.Printf("json umarshal failed, err: %+v", err)
				continue
			}
			data.Meta.ReceiveTime = time.Now()
			handler(&data)
		case <-ctx.Done():
			log.Println("context done")
			return nil
		}
	}
}

func (r *RedisSource) mockMessage(topic string) {
	for {
		marshal, _ := json.Marshal(map[string]any{
			"topic":   "topic1",
			"type":    "test",
			"payload": `{"username": "123", "nickname": "123"}`,
		})
		r.redis.Publish(context.Background(), topic, marshal)
	}
}
