package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"orchestrator/internal/config"
	"orchestrator/internal/logger"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
)

var (
	RestartDelaySeconds = 10
)

type OrchestratorInterface interface {
	ConsumeNotifications()

	Set(ctx context.Context) error
	Unset(ctx context.Context) error
	Status(ctx context.Context) ([]Status, error)

	RunAll(ctx context.Context) error
	RunGroup(ctx context.Context, group int) error
	Run(ctx context.Context, processUUID uuid.UUID) error

	StopAll(ctx context.Context) error
	StopGroup(ctx context.Context, group int) error
	Stop(ctx context.Context, processUUID uuid.UUID) error

	ExecLogs(ctx context.Context, logsType string, processUUID uuid.UUID, offset int) (string, error)
}

type Orchestrator struct {
	Logger        *log.Logger
	LoggerCleanup func()
	Notifications chan Notification
	Executables   Executables
}

type Notification struct {
	Executable *Executable
	err        error
}

func NewOrchestrator() *Orchestrator {
	logger, cleanup := logger.NewLogger()

	return &Orchestrator{
		Logger:        logger,
		LoggerCleanup: cleanup,
		Notifications: make(chan Notification),
		Executables:   make(Executables, 0),
	}
}

func (o *Orchestrator) ConsumeNotifications() {
	defer func() {
		close(o.Notifications)
		o.Logger.Print(logger.LogInfo + "Stopped consuming notifications")
	}()

	o.Logger.Print(logger.LogInfo + "Starting consuming Orchestrator notifications...")
	for notification := range o.Notifications {
		executable := notification.Executable

		if notification.err != nil {
			o.Logger.Printf(logger.LogErr+"Executable %s finished with error: %s", executable.Name, notification.err.Error())
		} else {
			o.Logger.Printf(logger.LogInfo+"Executable %s has finished successfully", executable.Name)
		}

		if executable.AutoRestart && !o.isErrorGracefull(notification.err) {
			o.Logger.Printf(logger.LogInfo+"Sleepin delay before starting the executable: %s", executable.Name)
			time.Sleep(time.Duration(RestartDelaySeconds) * time.Second)
			o.startExecutable(executable)
		}
	}
}

func (o *Orchestrator) Set(ctx context.Context) error {
	var err error

	if len(o.Executables) > 0 {
		return errors.New("executables already set")
	}

	file, err := os.Open(config.GetConfig().EXECUTABLES_JSON_PATH)
	if err != nil {
		return errors.New("error opening executables file: " + err.Error())
	}
	defer file.Close()

	var executables Executables
	err = json.NewDecoder(file).Decode(&executables)
	if err != nil {
		return errors.New("error decoding executables file: " + err.Error())
	}

	for _, executable := range executables {
		err = executable.validate()
		if err != nil {
			return errors.New("error validating executable: " + executable.Name + " " + err.Error())
		}
	}

	for _, executable := range executables {
		executable.ID = uuid.New()
	}

	o.Executables = executables

	return nil
}

func (o *Orchestrator) Unset(ctx context.Context) error {
	if len(o.Executables) == 0 {
		return errors.New("no executables to unset")
	}

	for _, executable := range o.Executables {
		if executable.status().Running {
			return errors.New("cannot unset executable " + executable.Name + " is running")
		}
	}

	o.Executables = make(Executables, 0)

	return nil
}

func (o *Orchestrator) Status(ctx context.Context) ([]Status, error) {
	statuses := make([]Status, 0, len(o.Executables))
	for _, executable := range o.Executables {
		status := executable.status()
		statuses = append(statuses, status)
	}

	return statuses, nil
}

/*
Current strategy: Start as many executables as possible. If an executable fails to start, log the error and continue.
Consider changing the strategy to force start all executables. If an executable fails to start, log the error and stop all executables.
*/
func (o *Orchestrator) RunAll(ctx context.Context) error {
	if len(o.Executables) == 0 {
		return errors.New("there are no executables set to run")
	}

	for _, executable := range o.Executables {
		o.startExecutable(executable)
	}

	return nil
}

func (o *Orchestrator) RunGroup(ctx context.Context, group int) error {
	executablesGroup := Executables{}

	for _, executable := range o.Executables {
		if executable.Group == group {
			executablesGroup = append(executablesGroup, executable)
		}
	}

	if len(executablesGroup) == 0 {
		return errors.New("no executables found in group")
	}

	for _, executable := range executablesGroup {
		o.startExecutable(executable)
	}

	return nil
}

func (o *Orchestrator) Run(ctx context.Context, processUUID uuid.UUID) error {
	var executable *Executable

	for _, exec := range o.Executables {
		if exec.ID == processUUID {
			executable = exec
			break
		}
	}

	if executable == nil {
		return errors.New("executable not found")
	}

	o.startExecutable(executable)

	return nil

}

func (o *Orchestrator) StopAll(ctx context.Context) error {
	if len(o.Executables) == 0 {
		return errors.New("no executables to stop")
	}

	for _, executable := range o.Executables {
		err := executable.stop()
		if err != nil {
			o.Logger.Printf(logger.LogErr+"Error stopping executable %s: %s", executable.Name, err.Error())
		}
	}

	return nil
}

func (o *Orchestrator) StopGroup(ctx context.Context, group int) error {
	executablesGroup := Executables{}

	for _, executable := range o.Executables {
		if executable.Group == group {
			executablesGroup = append(executablesGroup, executable)
		}
	}

	if len(executablesGroup) == 0 {
		return errors.New("no executables found in group")
	}

	for _, executable := range executablesGroup {
		err := executable.stop()
		if err != nil {
			o.Logger.Printf(logger.LogErr+"Error stopping executable %s: %s", executable.Name, err.Error())
		}
	}

	return nil
}

func (o *Orchestrator) Stop(ctx context.Context, processUUID uuid.UUID) error {
	var executable *Executable

	for _, exec := range o.Executables {
		if exec.ID == processUUID {
			executable = exec
			break
		}
	}

	if executable == nil {
		return errors.New("executable not found")
	}

	err := executable.stop()
	if err != nil {
		return errors.New("error stopping executable " + executable.Name + ": " + err.Error())
	}

	return nil
}

func (o *Orchestrator) ExecLogs(ctx context.Context, logsType string, processUUID uuid.UUID, offset int) (string, error) {
	var executable *Executable

	for _, exec := range o.Executables {
		if exec.ID == processUUID {
			executable = exec
			break
		}
	}

	if executable == nil {
		return "", errors.New("executable not found")
	}

	var logPrefix string
	switch logsType {
	case logger.LogTypeOut:
		logPrefix = executable.LogFileName
	case logger.LogTypeError:
		logPrefix = executable.ErrorFileName
	default:
		return "", errors.New("invalid logs type")
	}

	logsFolder := executable.LogDir
	files, err := os.ReadDir(logsFolder)
	if err != nil {
		return "", errors.New("error reading logs folder: " + err.Error())
	}

	var logs []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".log" && strings.HasPrefix(file.Name(), logPrefix) {
			logs = append(logs, file.Name())
		}
	}

	if len(logs) == 0 {
		return "", errors.New("no logs found")
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i] > logs[j]
	})

	if offset >= len(logs) {
		return "", errors.New("offset out of range")
	}

	offsetLog := logs[offset]
	filePath := filepath.Join(logsFolder, offsetLog)

	logContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.New("error opening log file: " + err.Error())
	}

	return string(logContent), nil
}

func (o *Orchestrator) startExecutable(executable *Executable) {
	if executable.status().Running {
		o.Logger.Printf(logger.LogInfo+"Executable %s is already running", executable.Name)
		return
	}

	err := executable.start()
	if err != nil {
		o.Logger.Printf(logger.LogErr+"Error on trying to start the executable %s : %s", executable.Name, err.Error())
		return
	}
	o.Logger.Printf(logger.LogInfo+"Executable %s started successfully", executable.Name)

	go executable.wait(o.Notifications)
}

func (o *Orchestrator) isErrorGracefull(err error) bool {
	if err == nil {
		return true
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		status := exitErr.Sys().(syscall.WaitStatus)
		if status.Signaled() && status.Signal() == GracefullExitSignal {
			return true
		}
	}

	return false
}
