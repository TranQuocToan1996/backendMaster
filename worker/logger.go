package worker

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerAsynq struct{}

func NewLoggerAsynq() *LoggerAsynq {
	return &LoggerAsynq{}
}

func (logger *LoggerAsynq) Print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (logger *LoggerAsynq) Printf(ctx context.Context, format string, v ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (logger *LoggerAsynq) Debug(args ...interface{}) {
	logger.Print(zerolog.DebugLevel, args...)
}

func (logger *LoggerAsynq) Info(args ...interface{}) {
	logger.Print(zerolog.InfoLevel, args...)
}

func (logger *LoggerAsynq) Warn(args ...interface{}) {
	logger.Print(zerolog.WarnLevel, args...)
}

func (logger *LoggerAsynq) Error(args ...interface{}) {
	logger.Print(zerolog.ErrorLevel, args...)
}

func (logger *LoggerAsynq) Fatal(args ...interface{}) {
	logger.Print(zerolog.FatalLevel, args...)
}
