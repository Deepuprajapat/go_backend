package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	logger *zerolog.Logger
)

func Set(l *zerolog.Logger) {
	logger = l
}

func Get() *zerolog.Logger {
	return logger
}

func Init() {
	// Set time format to Unix timestamp
	zerolog.TimeFieldFormat = time.RFC3339

	// Set global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("ENV") == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Create console writer with better formatting
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	// Create logger with console writer
	l := zerolog.New(output).With().Timestamp().Caller().Logger()
	logger = &l

	// Set as default logger
	Set(logger)
}
