package main

import (
	"context"
	"log"
	"orchestrator/internal/config"
	"orchestrator/internal/orchestrator"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	c := config.GetConfig()

	instance := orchestrator.NewOrchestrator()

	defer func() {
		instance.LoggerCleanup()
	}()

	e := echo.New()
	router := dependencies(instance)
	router.Route(e)

	go func() {
		e.Logger.Fatal(e.Start(":" + c.HTTP_PORT))
	}()

	time.Sleep(1 * time.Second)

	go instance.ConsumeNotifications()

	time.Sleep(1 * time.Second)

	if c.AUTOSETRUN {
		var err error
		ctx := context.Background()

		err = instance.Set(ctx)
		if err != nil {
			log.Fatal("Failed to set orchestrator: " + err.Error())
		}

		err = instance.RunAll(ctx)
		if err != nil {
			log.Fatal("Failed to run all processes: " + err.Error())
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")
}
