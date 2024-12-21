package logger

import "log"

var (
	LogInfo = "INFO: "
	LogErr  = "ERROR: "

	LoggingTimestampFormat = "020106-1504"
)

func NewLogger() *log.Logger {
	return log.Default()
}
