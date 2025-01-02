package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	LogInfo                = "INFO: "
	LogErr                 = "ERROR: "
	LoggingTimestampFormat = "06-01-02-15-04-05"
	OrchestratorLogs       = "/logs/"
	LogTypeError           = "errors"
	LogTypeOut             = "out"
)

func NewLogger() (*log.Logger, func()) {
	cwd, _ := os.Getwd()

	timestamp := time.Now().Format(LoggingTimestampFormat)
	logFilePath := cwd + OrchestratorLogs + fmt.Sprintf("%s-%s.log", "orchestrator", timestamp)
	file, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	mw := io.MultiWriter(os.Stdout, file)

	logger := log.New(mw, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)

	cleanup := func() {
		_ = file.Close()
	}

	return logger, cleanup
}
