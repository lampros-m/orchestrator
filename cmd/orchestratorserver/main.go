package main

import (
	"orchestrator/internal/config"

	"github.com/labstack/echo/v4"
)

func main() {
	c := config.GetConfig()

	router := dependencies()

	e := echo.New()
	router.Route(e)

	e.Logger.Fatal(e.Start(":" + c.HTTP_PORT))
}
