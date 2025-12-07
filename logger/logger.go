package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Level int

const (
	DebugLevel Level = -4
	InfoLevel  Level = 0
	WarnLevel  Level = 4
	ErrorLevel Level = 8
)

func SetupLogging(
	consoleLogLevel Level,
	fileLogLevel Level,
	logFilePath string,
) error {

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open log file: %v", err)
	}

	multi := newMultiHandler(
		newCustomConsoleHandler(
			slog.Level(consoleLogLevel),
		),
		slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level:     slog.Level(fileLogLevel),
			AddSource: true,
		}),
	)

	slog.SetDefault(slog.New(multi))
	return nil
}

func Debug(format string, args ...any) {
	DebugWithArgs(fmt.Sprintf(format, args...))
}

func DebugWithArgs(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(format string, args ...any) {
	InfoWithArgs(fmt.Sprintf(format, args...))
}

func InfoWithArgs(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(format string, args ...any) {
	WarnWithArgs(fmt.Sprintf(format, args...))
}

func WarnWithArgs(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(format string, args ...any) {
	ErrorWithArgs(fmt.Sprintf(format, args...))
}

func ErrorWithArgs(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Fatal(err error, args ...any) {
	slog.Error(fmt.Sprintf("Fatal Error: %v", err), args...)
	os.Exit(1)
}
