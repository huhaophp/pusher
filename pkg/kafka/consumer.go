package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"pusher/config"
)

func InitConsumer(conf *config.KafkaConfig) (*kafka.Consumer, error) {
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": conf.Servers,
		"group.id":          "group2",
		"auto.offset.reset": "earliest",
	})
}
