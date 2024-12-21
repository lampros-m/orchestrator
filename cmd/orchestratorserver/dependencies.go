package main

import (
	"orchestrator/internal/apihttp"
	"orchestrator/internal/apihttp/controllers"
	"orchestrator/internal/orchestrator"
)

func dependencies() *apihttp.Router {
	// orchestrator
	instance := orchestrator.NewOrchestrator()

	// controllers
	orchestrator := controllers.NewOrchestrator(instance)

	return apihttp.NewRouter(
		orchestrator,
	)
}
