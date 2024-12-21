package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"orchestrator/internal/config"
	"orchestrator/internal/logger"
	"os"
	"sync"

	"github.com/google/uuid"
)

type OrchestratorInterface interface {
	Set(ctx context.Context) error
	Unset(ctx context.Context) error
	Run(ctx context.Context) error
	Status(ctx context.Context) ([]Status, error)
	StopAll(ctx context.Context) error
}

type Orchestrator struct {
	Logger        *log.Logger
	wg            *sync.WaitGroup
	Notifications chan Notification
	Executables   Executables
}

type Notification struct {
	Executable *Executable
	err        error
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		Logger:        logger.NewLogger(),
		wg:            &sync.WaitGroup{},
		Notifications: make(chan Notification),
		Executables:   make(Executables, 0),
	}
}

func (o *Orchestrator) Set(ctx context.Context) error {
	var err error

	if len(o.Executables) > 0 {
		o.Logger.Print(logger.LogErr + "executables already set")
		return errors.New("executables already set")
	}

	file, err := os.Open(config.GetConfig().EXECUTABLES_JSON_PATH)
	if err != nil {
		o.Logger.Print(logger.LogErr + "error opening executables file: " + err.Error())
		return errors.New("error opening executables file: " + err.Error())
	}
	defer file.Close()

	var executables Executables
	err = json.NewDecoder(file).Decode(&executables)
	if err != nil {
		o.Logger.Print(logger.LogErr + "error decoding executables file: " + err.Error())
		return errors.New("error decoding executables file: " + err.Error())
	}

	for _, executable := range executables {
		err = executable.validate()
		if err != nil {
			o.Logger.Printf(logger.LogErr+"error validating executable: %s %s", executable.Name, err.Error())
			return errors.New("error validating executable: " + executable.Name + " " + err.Error())
		}
	}

	for _, executable := range executables {
		executable.ID = uuid.New()
	}

	o.Executables = executables
	o.Logger.Print(logger.LogInfo + "set completed successfully")

	return nil
}

func (o *Orchestrator) Unset(ctx context.Context) error {
	if len(o.Executables) == 0 {
		o.Logger.Print(logger.LogErr + "no executables to unset")
		return errors.New("no executables to unset")
	}

	for _, executable := range o.Executables {
		if executable.Status().Running {
			o.Logger.Print(logger.LogErr + "cannot unset executable " + executable.Name + " is running")
			return errors.New("cannot unset executable " + executable.Name + " is running")
		}
	}

	o.Executables = make(Executables, 0)
	o.Logger.Print(logger.LogInfo + "unset completed successfully")

	return nil
}

/*
Current strategy: Start as many executables as possible. If an executable fails to start, log the error and continue.
Consider changing the strategy to force start all executables. If an executable fails to start, log the error and stop all executables.

Consider the possibility of dependencies between executables. If an executable depends on another, it should wait for the other to start before starting itself.
*/
func (o *Orchestrator) Run(ctx context.Context) error {
	if len(o.Executables) == 0 {
		o.Logger.Print(logger.LogErr + "there are no executables set to run")
		return errors.New("there are no executables set to run")
	}

	for _, executable := range o.Executables {
		o.startExecutable(executable)
	}

	go o.waitGroupWait()
	go o.consumeNotifications()
	o.Logger.Print(logger.LogInfo + "run completed successfully")

	return nil
}

func (o *Orchestrator) Status(ctx context.Context) ([]Status, error) {
	statuses := make([]Status, 0, len(o.Executables))
	for _, executable := range o.Executables {
		status := executable.Status()
		statuses = append(statuses, status)
	}

	o.Logger.Print(logger.LogInfo + "status completed successfully")

	return statuses, nil
}

/*
Stop all terminates all running processes. If an executable fails to stop, log the error and continue. It tries to terminate as many processes as possible.
But, if an executable is configured to restart automatically, it will be restarted after stopping.
Consider changing the strategy to force stop all executables and not restart any of them - even if they are configured to restart automatically.
*/
func (o *Orchestrator) StopAll(ctx context.Context) error {
	if len(o.Executables) == 0 {
		o.Logger.Print(logger.LogErr + "no executables to stop")
		return errors.New("no executables to stop")
	}

	for _, executable := range o.Executables {
		err := executable.Stop()
		if err != nil {
			o.Logger.Printf("error stopping executable %s: %s", executable.Name, err.Error())
		}
	}

	o.Logger.Print(logger.LogInfo + "stop all completed successfully")

	return nil
}

func (o *Orchestrator) startExecutable(executable *Executable) {
	if executable.Status().Running {
		o.Logger.Printf(logger.LogInfo+"executable %s is already running", executable.Name)
		return
	}

	err := executable.Start()
	if err != nil {
		o.Logger.Printf(logger.LogErr+"error on trying to start the executable %s : %s", executable.Name, err.Error())
		return
	}
	o.Logger.Printf(logger.LogInfo+"executable %s started successfully", executable.Name)
	o.wg.Add(1)
	go executable.Wait(o.wg, o.Notifications)
}

func (o *Orchestrator) consumeNotifications() {
	o.Logger.Print(logger.LogInfo + "started consuming notifications")
	defer close(o.Notifications)

	for notification := range o.Notifications {
		executable := notification.Executable

		if notification.err != nil {
			o.Logger.Printf(logger.LogErr+"executable %s finished with error: %s", executable.Name, notification.err.Error())
		} else {
			o.Logger.Printf(logger.LogInfo+"executable %s has finished successfully", executable.Name)
		}

		if executable.AutoRestart {
			o.Logger.Printf(logger.LogInfo+"restarting executable %s", executable.Name)
			o.startExecutable(executable)
		}
	}

	o.Logger.Print(logger.LogInfo + "stopped consuming notifications")
}

func (o *Orchestrator) waitGroupWait() {
	o.Logger.Print(logger.LogInfo + "started waiting for executables")

	defer func() {
		// Placeholder for code when all executables are done
	}()

	o.wg.Wait()

	o.Logger.Print(logger.LogInfo + "stopped waiting for executables")
}
