package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)


func main() {
	_ , cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	Execute()
}