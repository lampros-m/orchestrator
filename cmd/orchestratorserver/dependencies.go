package main

import (
	"orchestrator/internal/apihttp"
	"orchestrator/internal/apihttp/controllers"
	"orchestrator/internal/orchestrator"
)

func dependencies(instance *orchestrator.Orchestrator) *apihttp.Router {
	orchestrator := controllers.NewOrchestrator(instance)

	return apihttp.NewRouter(
		orchestrator,
	)
}
