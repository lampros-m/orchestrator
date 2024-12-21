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
	e.GET("/set", o.Orchestrator.Set)
	e.GET("/unset", o.Orchestrator.Unset)
	e.GET("/status", o.Orchestrator.Status)
	e.GET("/run", o.Orchestrator.Run)
	e.GET("/stopall", o.Orchestrator.StopAll)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
