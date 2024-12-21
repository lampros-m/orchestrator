package main

import (
	"context"
	"log"
	"orchestrator/internal/orchestrator"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	instance := orchestrator.NewOrchestrator()

	ctx := context.Background()
	var err error

	err = instance.Set(ctx)
	if err != nil {
		log.Fatal("Failed to set orchestrator: " + err.Error())
	}

	err = instance.Run(ctx)
	if err != nil {
		log.Fatal("Failed to run orchestrator: " + err.Error())
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")
}
