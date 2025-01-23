package utils

import (
	"os"
	zerolog "github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger() *Logger {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	return &Logger{
		logger: logger,
	}
}

func (this_ *Logger) Error() *zerolog.Event {
   return this_.logger.Error()
}

func (this_ *Logger) Info() *zerolog.Event {
   return this_.logger.Info()
}

func (this_ *Logger) Debug() *zerolog.Event {
	return this_.logger.Debug()
}

func (this_ *Logger) Warn() *zerolog.Event {
	return this_.logger.Warn()
}

func (logger *Logger) Trace() *zerolog.Event {
	return logger.logger.Trace()
}

func (logger *Logger) Fatal() *zerolog.Event {
	return logger.logger.Fatal()
}

func (logger *Logger) Panic() *zerolog.Event {
	return logger.logger.Panic()
}
