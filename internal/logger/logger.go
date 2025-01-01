package logger

import (
	"context"
	"fmt"
)

type ILogger interface {
	Trace(context.Context, string)
	Debug(context.Context, string)
	Info(context.Context, string)
	Warn(context.Context, string)
	Error(context.Context, string)
}

type LogLevel int

const (
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var logger ILogger

func Init(log ILogger) {
	if logger == nil {
		logger = log
	} else {
		Error(context.Background(), "Logger was already initialized. Cannot initialize again!")
	}
}

func Debug(ctx context.Context, message string) {
	if logger != nil {
		logger.Debug(ctx, message)
	}
}

func DebugArgs(ctx context.Context, message string, args ...any) {
	Debug(ctx, fmt.Sprintf(message, args...))
}

func Trace(ctx context.Context, message string) {
	if logger != nil {
		logger.Trace(ctx, message)
	}
}

func TraceArgs(ctx context.Context, message string, args ...any) {
	Trace(ctx, fmt.Sprintf(message, args...))
}

func Info(ctx context.Context, message string) {
	if logger != nil {
		logger.Info(ctx, message)
	}
}

func InfoArgs(ctx context.Context, message string, args ...any) {
	Info(ctx, fmt.Sprintf(message, args...))
}

func Warn(ctx context.Context, message string) {
	if logger != nil {
		logger.Warn(ctx, message)
	}
}

func WarnArgs(ctx context.Context, message string, args ...any) {
	Warn(ctx, fmt.Sprintf(message, args...))
}

func Error(ctx context.Context, message string) {
	if logger != nil {
		logger.Error(ctx, message)
	}
}

func ErrorArgs(ctx context.Context, message string, args ...any) {
	Error(ctx, fmt.Sprintf(message, args...))
}
