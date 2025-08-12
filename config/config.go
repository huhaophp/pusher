package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APP       APP             `yaml:"app"`
	Redis     RedisConfig     `yaml:"redis"`
	Kafka     KafkaConfig     `yaml:"kafka"`
	Source    SourceConfig    `yaml:"source"`
	Logger    LoggerConfig    `yaml:"logger"`
	PProf     PProfConfig     `yaml:"pprof"`
	WebSocket WebSocketConfig `yaml:"webSocket"`
}

type APP struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type SourceConfig struct {
	Redis []string `yaml:"redis"`
	Kafka []string `yaml:"kafka"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type KafkaConfig struct {
	Servers         string `yaml:"servers"`
	GroupID         string `yaml:"groupID"`
	AutoOffsetReset string `yaml:"autoOffsetReset"`
}

type PProfConfig struct {
	Port   string `yaml:"port"`
	Enable bool   `yaml:"enable"`
}

type WebSocketConfig struct {
	SubscriberChanSize int `yaml:"subscriberChanSize"`
}

func LoadConfig(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, err
	}

	config.Kafka.GroupID = fmt.Sprintf("%s-%d", config.Kafka.GroupID, os.Getpid())

	return &config, nil
}
