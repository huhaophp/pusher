package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APP    APP          `yaml:"app"`
	Redis  RedisConfig  `yaml:"redis"`
	Source SourceConfig `yaml:"source"`
	Logger LoggerConfig `yaml:"logger"`
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

func LoadConfig(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
