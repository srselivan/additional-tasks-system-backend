package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	logsPath    = "_logs"
	serviceName = "backend"
)

// New creates new logger instance with specified logging level.
// Logger writes logs into stdout and file.
func New(level string) *zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(err)
	}

	output := zerolog.ConsoleWriter{
		TimeFormat: time.RFC3339Nano,
		Out:        os.Stdout,
	}

	if _, err = os.Stat(logsPath); os.IsNotExist(err) {
		if err = os.Mkdir(logsPath, 0777); err != nil {

			panic(err)
		}
	}

	logsFilePath := fmt.Sprintf("%s/%s.log", logsPath, serviceName)
	file, err := os.OpenFile(logsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}

	multi := zerolog.MultiLevelWriter(output, file)

	l := zerolog.New(multi).With().Caller().Timestamp().Logger()
	l.Level(logLevel)

	return &l
}
