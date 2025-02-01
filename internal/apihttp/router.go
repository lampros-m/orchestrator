package apihttp

import (
	_ "orchestrator/docs"
	"orchestrator/internal/apihttp/controllers"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	Orchestrator controllers.OrchestratorInterface
}

func NewRouter(
	orchestrator controllers.OrchestratorInterface,
) *Router {
	return &Router{
		Orchestrator: orchestrator,
	}
}

//	@title			orchestrator-api
//	@version		0.0.1
//	@description	This is an API that controls running processes.

// @BasePath
func (o *Router) Route(e *echo.Echo) {
	// Generic
	e.GET("/set", o.Orchestrator.Set)
	e.GET("/unset", o.Orchestrator.Unset)
	e.GET("/status", o.Orchestrator.Status)

	// Run
	e.GET("/runall", o.Orchestrator.RunAll)
	e.GET("/rungroup", o.Orchestrator.RunGroup)
	e.GET("/run", o.Orchestrator.Run)

	// Stop
	e.GET("/stopall", o.Orchestrator.StopAll)
	e.GET("/stopgroup", o.Orchestrator.StopGroup)
	e.GET("/stop", o.Orchestrator.Stop)

	// Logs
	e.GET("/execlogs", o.Orchestrator.ExecLogs)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Static
	e.Static("/", "assets")
	e.GET("index", func(c echo.Context) error { return c.File("assets/index.html") })
}
