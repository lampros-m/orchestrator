package orchestrator

import (
	"errors"
	"fmt"
	"io"
	"orchestrator/internal/logger"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type Executables []*Executable

type Executable struct {
	Configuration
	Process
}

type Configuration struct {
	Name          string   `json:"name"`
	BinaryPath    string   `json:"binary_path"`
	WorkingDir    string   `json:"working_dir"`
	LogDir        string   `json:"log_dir"`
	Arguments     []string `json:"arguments"`
	LogFileName   string   `json:"log_file_name"`
	ErrorFileName string   `json:"error_file_name"`
	AutoRestart   bool     `json:"auto_restart"`
}

type Process struct {
	ID                  uuid.UUID
	PID                 int
	CMD                 *exec.Cmd
	OutLogFileHandle    *os.File
	ErrorsLogFileHandle *os.File
}

type Status struct {
	ID          string `json:"id"`
	Running     bool   `json:"running"`
	Name        string `json:"name"`
	PID         int    `json:"pid"`
	AutoRestart bool   `json:"auto_restart"`
}

func (o *Executable) Start() error {
	if o.Status().Running {
		return errors.New("executable is already running: " + o.Name)
	}

	var err error

	timestamp := time.Now().Format(logger.LoggingTimestampFormat)

	logFilePath := o.LogDir + "/" + fmt.Sprintf("%s-%s.log", o.LogFileName, timestamp)
	errFilePath := o.LogDir + "/" + fmt.Sprintf("%s-%s.log", o.ErrorFileName, timestamp)

	var outLogF, errLogF *os.File
	defer func() {
		if err != nil {
			if outLogF != nil {
				outLogF.Close()
			}
			if errLogF != nil {
				errLogF.Close()
			}
		}
	}()

	outLogF, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file for %s: %w", o.Name, err)
	}

	errLogF, err = os.OpenFile(errFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open error file for %s: %w", o.Name, err)
	}

	cmd := exec.Command(o.BinaryPath, o.Arguments...)
	cmd.Dir = o.WorkingDir
	cmd.Stdout = io.MultiWriter(outLogF) // Place os.Stdout to see the output of each executable in the terminal - "dirty way"
	cmd.Stderr = io.MultiWriter(errLogF) // Place os.Stdout to see the output of each executable in the terminal - "dirty way"

	err = cmd.Start()
	if err != nil {
		return errors.New("error running executable command: " + o.Name + err.Error())
	}

	o.CMD = cmd
	o.OutLogFileHandle = outLogF
	o.ErrorsLogFileHandle = errLogF

	o.PID = cmd.Process.Pid

	return nil
}

func (o *Executable) Wait(wg *sync.WaitGroup, stopNotifications chan Notification) {
	defer wg.Done()

	err := o.CMD.Wait()

	o.OutLogFileHandle.Close()
	o.ErrorsLogFileHandle.Close()
	o.CMD = nil
	stopNotifications <- Notification{Executable: o, err: err}
}

func (o *Executable) Status() Status {
	status := Status{}
	status.ID = o.ID.String()
	status.Name = o.Name
	status.PID = o.PID
	status.AutoRestart = o.AutoRestart

	running := true
	switch {
	case o.CMD == nil, o.CMD.Process == nil:
		running = false
	default:
		running = o.CMD.Process.Signal(syscall.Signal(0)) == nil
	}
	status.Running = running

	return status
}

func (o *Executable) Stop() error {
	if o.CMD == nil || o.CMD.Process == nil {
		return nil
	}

	// Other options: SIGINT for interruption or Process.Kill() for imidiate termination
	if err := o.CMD.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to signal executable %s : %w", o.Name, err)
	}

	return nil
}

func (o *Executable) validate() error {
	// Name
	if o.Name == "" {
		return errors.New("executable name is required: " + o.Name)
	}

	// Binary Path
	if o.BinaryPath == "" {
		return errors.New("binary path is required: " + o.Name)
	}

	binaryPathInfo, err := os.Stat(o.BinaryPath)
	if err != nil {
		return errors.New("error stating binary path: " + o.Name)
	}
	if binaryPathInfo.IsDir() {
		return errors.New("binary path is a directory: " + o.Name)
	}

	if binaryPathInfo.Mode()&0111 == 0 {
		return errors.New("binary path is not executable: " + o.Name)
	}

	// Working Directory
	if o.WorkingDir == "" {
		return errors.New("executable working directory is required: " + o.Name)
	}

	workingDirectoryInfo, err := os.Stat(o.WorkingDir)
	if err != nil {
		return errors.New("error stating working directory: " + o.Name)
	}
	if !workingDirectoryInfo.IsDir() {
		return errors.New("working directory is not a directory of service: " + o.Name)
	}

	testFile := o.WorkingDir + "/.testwrite"
	f, err := os.Create(testFile)
	if err != nil {
		return errors.New("cannot write to working directory: " + o.Name)
	}
	f.Close()
	os.Remove(testFile)

	// Log Directory
	if o.LogDir == "" {
		return errors.New("log directory is required: " + o.Name)
	}

	logDirectoryInfo, err := os.Stat(o.LogDir)
	if err != nil {
		return errors.New("error stating log directory: " + o.Name)
	}
	if !logDirectoryInfo.IsDir() {
		return errors.New("log directory is not a directory: " + o.Name)
	}

	testFile = o.LogDir + "/.testwrite"
	f, err = os.Create(testFile)
	if err != nil {
		return errors.New("cannot write to log directory: " + o.Name)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}
