package logger

import (
	"os"
	"pusher/config"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

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

func Info(args ...any) {
	log.Info(args...)
}

func Infof(format string, args ...any) {
	log.Infof(format, args...)
}

func Error(args ...any) {
	log.Error(args...)
}

func Errorf(format string, args ...any) {
	log.Errorf(format, args...)
}

func Warn(args ...any) {
	log.Warn(args...)
}

func Warnf(format string, args ...any) {
	log.Warnf(format, args...)
}

func Debug(args ...any) {
	log.Debug(args...)
}

func Debugf(format string, args ...any) {
	log.Debugf(format, args...)
}

func Fatal(args ...any) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	log.Fatalf(format, args...)
}
