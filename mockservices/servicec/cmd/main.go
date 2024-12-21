package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	ServiceName = "Service Charlie"
)

func main() {
	fmt.Println("Hello, World! I'm service: " + ServiceName)

	time.Sleep(45 * time.Second)

	log.Println("Service is shutting down...")
	os.Exit(130)
}
