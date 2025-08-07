package source

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"pusher/internal/types"
	"pusher/pkg/logger"
	"time"
)

type KafkaSource struct {
	consumer *kafka.Consumer
}

func NewKafkaSource(consumer *kafka.Consumer) *KafkaSource {
	return &KafkaSource{consumer: consumer}
}

func (k *KafkaSource) PullMessage(ctx context.Context, topic string, handler func(data *types.Data)) error {
	err := k.consumer.SubscribeTopics([]string{topic}, func(c *kafka.Consumer, e kafka.Event) error {
		switch ev := e.(type) {
		case kafka.AssignedPartitions:
			logger.Infof("Assigned partitions: %+v", ev.Partitions)
			return c.Assign(ev.Partitions)
		case kafka.RevokedPartitions:
			logger.Infof("Revoked partitions: %+v", ev.Partitions)
			return c.Unassign()
		}
		return nil
	})
	if err != nil {
		return err
	}

	defer k.consumer.Close()

	for {
		select {
		case <-ctx.Done():
			logger.Infof("topic: %s context done.", topic)
			return nil
		default:
			msg, err := k.consumer.ReadMessage(time.Second)
			if err == nil {
				var data types.Data
				if err := json.Unmarshal(msg.Value, &data); err != nil {
					logger.Infof("json umarshal failed, err: %+v", err)
					continue
				}
				data.Meta.ReceiveTime = time.Now()
				handler(&data)
			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}
}
