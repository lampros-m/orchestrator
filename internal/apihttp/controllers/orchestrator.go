package controllers

import (
	"net/http"
	"orchestrator/internal/apihttp/dtos"
	"orchestrator/internal/orchestrator"

	"github.com/labstack/echo/v4"
)

type OrchestratorInterface interface {
	Set(c echo.Context) error
	Unset(c echo.Context) error
	Run(c echo.Context) error
	Status(c echo.Context) error
	StopAll(echoContext echo.Context) error
}

type Orchestrator struct {
	instance orchestrator.OrchestratorInterface
}

func NewOrchestrator(
	instance orchestrator.OrchestratorInterface,
) *Orchestrator {
	return &Orchestrator{
		instance: instance,
	}
}

// Set godoc
//
//	@Summary		Set the executables
//	@Description	This endpoint sets the executables in the orchestrator. In order to set the executables again, all processes must be stopped and unset.
//	@Tags			orchestrator
//	@Produce		json
//	@Success		200	{object}	dtos.SetResponse
//	@Failure		500	{object}	dtos.SetResponse
//	@Router			/set [get]
func (o *Orchestrator) Set(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	err := o.instance.Set(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.SetResponse{Message: "Failed to set orchestrator: " + err.Error()})
	}

	response := dtos.SetResponse{Message: "Orchestrator set successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// Unset godoc
//
//	@Summary		Unset the executables
//	@Description	This endpoint unsets the executables in the orchestrator. In order to unset the executables, all processes must be stopped.
//	@Tags			orchestrator
//	@Produce		json
//	@Success		200	{object}	dtos.UnsetResponse
//	@Failure		500	{object}	dtos.UnsetResponse
//	@Router			/unset [get]
func (o *Orchestrator) Unset(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	err := o.instance.Unset(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.UnsetResponse{Message: "Failed to unset orchestrator: " + err.Error()})
	}

	response := dtos.UnsetResponse{Message: "Orchestrator unset successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// Run godoc
//
//	@Summary		Run all the executables
//	@Description	This endpoint tries to run all the executables that are set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		json
//	@Success		200	{object}	dtos.RunResponse
//	@Failure		500	{object}	dtos.RunResponse
//	@Router			/run [get]
func (o *Orchestrator) Run(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	err := o.instance.Run(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.RunResponse{Message: "Failed to run orchestrator: " + err.Error()})
	}

	response := dtos.RunResponse{Message: "Orchestrator started successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// Status godoc
//
//	@Summary		Return the status of the executables
//	@Description	This endpoint returns the status of the executables that are set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		json
//	@Success		200	{object}	[]orchestrator.Status
//	@Failure		500	{object}	dtos.StatusResponse
//	@Router			/status [get]
func (o *Orchestrator) Status(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	status, err := o.instance.Status(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.StatusResponse{Message: "Failed to get orchestrator status: " + err.Error()})
	}

	return echoContext.JSON(http.StatusOK, status)
}

// StopAll godoc
//
//	@Summary		Stops all the executables
//	@Description	This endpoint tries to stop all the executables that are set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		json
//	@Success		200	{object}	dtos.StopAllResponse
//	@Failure		500	{object}	dtos.StopAllResponse
//	@Router			/stopall [get]
func (o *Orchestrator) StopAll(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	err := o.instance.StopAll(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.StopAllResponse{Message: "Failed to stop all orchestrator: " + err.Error()})
	}

	response := dtos.StopAllResponse{Message: "Orchestrator stopped successfully"}

	return echoContext.JSON(http.StatusOK, response)
}
