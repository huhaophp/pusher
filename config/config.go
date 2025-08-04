package config

type Config struct {
	APP   APP         `yaml:"app"`
	Redis RedisConfig `yaml:"redis"`
}

type APP struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Topic struct {
	Name string `yaml:"name"`
}
