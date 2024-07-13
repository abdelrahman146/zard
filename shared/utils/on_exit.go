package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func OnExit(cbs ...func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, f := range cbs {
		f()
	}
}
