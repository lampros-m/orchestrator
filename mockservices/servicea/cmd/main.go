package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	ServiceName = "Service Alpha"
)

func main() {
	fmt.Println("Hello, World! I'm service: " + ServiceName)

	time.Sleep(15 * time.Second)

	log.Println("Service is shutting down...")
	os.Exit(1)
}
