package source

import (
	"context"
	"encoding/json"
	"errors"
	"pusher/config"
	"pusher/internal/types"
	"pusher/pkg/logger"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaSource struct {
	config *config.KafkaConfig
}

func NewKafkaSource(config *config.KafkaConfig) *KafkaSource {
	return &KafkaSource{config: config}
}

func (k *KafkaSource) PullMessage(ctx context.Context, topic string, handler func(data *types.Data)) error {
	for {
		err := k.consumeLoop(ctx, topic, handler)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil // 上下文结束，正常退出
			}
			logger.GetLogger().Errorf("Kafka consume error: %v, retrying in 5s...", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func (k *KafkaSource) consumeLoop(ctx context.Context, topic string, handler func(data *types.Data)) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     strings.Split(k.config.Servers, ","),
		GroupID:     k.config.GroupID,
		StartOffset: kafka.LastOffset,
		Topic:       topic,
	})
	defer reader.Close()
	logger.GetLogger().WithField("topic", topic).Info("kafka consumer started")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				return err // 交给外层重连
			}

			var data types.Data
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				logger.GetLogger().Errorf("Unmarshal error: %v", err)
				continue
			}

			handler(&data)
		}
	}
}
