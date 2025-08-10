package logger

import (
	"os"
	"pusher/config"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func Init(conf *config.LoggerConfig) {
	log = logrus.New()

	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		level = logrus.InfoLevel
	}

	log.SetLevel(level)

	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetOutput(os.Stdout)
}

func GetLogger() *logrus.Logger {
	return log
}
