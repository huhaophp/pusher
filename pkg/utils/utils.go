package utils

import (
	"context"
	"os"
	"os/signal"
	"pusher/pkg/logger"
	"syscall"
	"time"
)

func WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Infof("received signal: %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	<-ctx.Done()
}
