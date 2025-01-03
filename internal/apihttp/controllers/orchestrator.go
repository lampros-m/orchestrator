package controllers

import (
	"net/http"
	"orchestrator/internal/apihttp/dtos"
	"orchestrator/internal/orchestrator"
	"strconv"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

type OrchestratorInterface interface {
	Set(echoContext echo.Context) error
	Unset(echoContext echo.Context) error
	Status(echoContext echo.Context) error
	RunAll(echoContext echo.Context) error
	RunGroup(echoContext echo.Context) error
	Run(echoContext echo.Context) error
	StopAll(echoContext echo.Context) error
	StopGroup(echoContext echo.Context) error
	Stop(echoContext echo.Context) error
	ExecLogs(echoContext echo.Context) error
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
//	@Success		200	{object}	dtos.GenericResponse
//	@Failure		500	{object}	dtos.GenericResponse
//	@Router			/set [get]
func (o *Orchestrator) Set(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	err := o.instance.Set(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.GenericResponse{Message: "Failed to set orchestrator: " + err.Error()})
	}

	response := dtos.GenericResponse{Message: "Orchestrator set successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// Unset godoc
//
//	@Summary		Unset the executables
//	@Description	This endpoint unsets the executables in the orchestrator. In order to unset the executables, all processes must be stopped.
//	@Tags			orchestrator
//	@Produce		json
//	@Success		200	{object}	dtos.GenericResponse
//	@Failure		500	{object}	dtos.GenericResponse
//	@Router			/unset [get]
func (o *Orchestrator) Unset(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	err := o.instance.Unset(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.GenericResponse{Message: "Failed to unset orchestrator: " + err.Error()})
	}

	response := dtos.GenericResponse{Message: "Orchestrator unset successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// Status godoc
//
//	@Summary		Return the status of the executables
//	@Description	This endpoint returns the status of the executables that are set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		json
//	@Success		200	{object}	[]orchestrator.Status
//	@Failure		500	{object}	dtos.GenericResponse
//	@Router			/status [get]
func (o *Orchestrator) Status(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	status, err := o.instance.Status(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.GenericResponse{Message: "Failed to get orchestrator status: " + err.Error()})
	}

	return echoContext.JSON(http.StatusOK, status)
}

// Run godoc
//
//	@Summary		Run all the executables
//	@Description	This endpoint tries to run all the executables that are set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		json
//	@Success		200	{object}	dtos.GenericResponse
//	@Failure		500	{object}	dtos.GenericResponse
//	@Router			/runall [get]
func (o *Orchestrator) RunAll(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	err := o.instance.RunAll(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.GenericResponse{Message: "Failed to run orchestrator: " + err.Error()})
	}

	response := dtos.GenericResponse{Message: "Orchestrator started successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// RunGroup godoc
//
//	@Summary		Run a group of executables
//	@Description	This endpoint tries to run a group of executables that are set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		json
//	@Param			group	query		int	true	"Group ID to run"
//	@Success		200		{object}	dtos.GenericResponse
//	@Failure		500		{object}	dtos.GenericResponse
//	@Router			/rungroup [get]
func (o *Orchestrator) RunGroup(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	group := echoContext.QueryParam("group")
	groupInt, _ := strconv.Atoi(group)
	if groupInt == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, dtos.GenericResponse{Message: "Group cannot be parsed successfully"})
	}

	err := o.instance.RunGroup(ctx, groupInt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.GenericResponse{Message: "Failed to run group: " + err.Error()})
	}

	response := dtos.GenericResponse{Message: "Group started successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// Run godoc
//
//	@Summary		Run an executable
//	@Description	This endpoint tries to run an executable that is set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		json
//	@Param			id	query		string	true	"UUID of the executable to run"	format(uuid)
//	@Success		200	{object}	dtos.GenericResponse
//	@Failure		500	{object}	dtos.GenericResponse
//	@Router			/run [get]
func (o *Orchestrator) Run(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	executableID := echoContext.QueryParam("id")
	executableUUID, err := uuid.Parse(executableID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dtos.GenericResponse{Message: "Cannot parse process ID as UUID: " + err.Error()})
	}

	err = o.instance.Run(ctx, executableUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.GenericResponse{Message: "Failed to run process: " + err.Error()})
	}

	response := dtos.GenericResponse{Message: "Process started successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// StopAll godoc
//
//	@Summary		Stops all the executables
//	@Description	This endpoint tries to stop all the executables that are set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		json
//	@Success		200	{object}	dtos.GenericResponse
//	@Failure		500	{object}	dtos.GenericResponse
//	@Router			/stopall [get]
func (o *Orchestrator) StopAll(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	err := o.instance.StopAll(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.GenericResponse{Message: "Failed to stop all orchestrator: " + err.Error()})
	}

	response := dtos.GenericResponse{Message: "Orchestrator stopped successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// StopGroup godoc
//
//	@Summary		Stops a group of executables
//	@Description	This endpoint tries to stop a group of executables that are set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		json
//	@Param			group	query		int	true	"Group ID to stop"
//	@Success		200		{object}	dtos.GenericResponse
//	@Failure		500		{object}	dtos.GenericResponse
//	@Router			/stopgroup [get]
func (o *Orchestrator) StopGroup(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	group := echoContext.QueryParam("group")
	groupInt, _ := strconv.Atoi(group)
	if groupInt == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, dtos.GenericResponse{Message: "Group cannot be parsed successfully"})
	}

	err := o.instance.StopGroup(ctx, groupInt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.GenericResponse{Message: "Failed to stop group: " + err.Error()})
	}

	response := dtos.GenericResponse{Message: "Group stopped successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// Stop godoc
//
//	@Summary		Stops an executable
//	@Description	This endpoint tries to stop an executable that is set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		json
//	@Param			id	query		string	true	"UUID of the group to stop"	format(uuid)
//	@Success		200	{object}	dtos.GenericResponse
//	@Failure		500	{object}	dtos.GenericResponse
//	@Router			/stop [get]
func (o *Orchestrator) Stop(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	executableID := echoContext.QueryParam("id")
	executableUUID, err := uuid.Parse(executableID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dtos.GenericResponse{Message: "Cannot parse process ID as UUID: " + err.Error()})
	}

	err = o.instance.Stop(ctx, executableUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.GenericResponse{Message: "Failed to stop process: " + err.Error()})
	}

	response := dtos.GenericResponse{Message: "Process stopped successfully"}

	return echoContext.JSON(http.StatusOK, response)
}

// ExecLogs godoc
//
//	@Summary		Get the logs of an executable
//	@Description	This endpoint tries to get the logs of an executable that is set in the orchestrator.
//	@Tags			orchestrator
//	@Produce		text/plain
//	@Param			id		query		string	true	"UUID of the executable to get logs"	format(uuid)
//	@Param			type	query		string	true	"Type of logs to get"					enum("errors", "out")
//	@Param			offset	query		int		false	"Offset of the logs to get"				default(0)
//	@Success		200		{string}	string
//	@Failure		500		{object}	dtos.GenericResponse
//	@Router			/execlogs [get]
func (o *Orchestrator) ExecLogs(echoContext echo.Context) error {
	ctx := echoContext.Request().Context()

	executableID := echoContext.QueryParam("id")
	executableUUID, err := uuid.Parse(executableID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dtos.GenericResponse{Message: "Cannot parse process ID as UUID: " + err.Error()})
	}

	logsType := echoContext.QueryParam("type")
	if logsType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, dtos.GenericResponse{Message: "Logs type is required"})
	}

	offset := echoContext.QueryParam("offset")
	offsetInt, _ := strconv.Atoi(offset)

	logs, err := o.instance.ExecLogs(ctx, logsType, executableUUID, offsetInt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dtos.GenericResponse{Message: "Failed to get logs: " + err.Error()})
	}

	return echoContext.String(http.StatusOK, logs)
}
