package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

// Init 初始化日志系统
func Init() {
	log = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel("debug")
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	log.SetFormatter(&logrus.JSONFormatter{})

	// 设置输出
	log.SetOutput(os.Stdout)
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	if log == nil {
		Init()
	}
	return log
}

// Info 记录信息日志
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof 记录格式化信息日志
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Error 记录错误日志
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf 记录格式化错误日志
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Warn 记录警告日志
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf 记录格式化警告日志
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Debug 记录调试日志
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf 记录格式化调试日志
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Fatal 记录致命错误日志并退出
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf 记录格式化致命错误日志并退出
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}
