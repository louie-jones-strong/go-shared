package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// Handler wrapper that combines the attributes into the final desired string
type customConsoleHandler struct {
	slog.Handler
}

// newCustomConsoleHandler creates a new CustomConsoleHandler instance
func newCustomConsoleHandler(
	logLevel slog.Level,
) *customConsoleHandler {
	return &customConsoleHandler{
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: false,
		}),
	}
}

func getColouredLevelLabel(level slog.Level) string {

	switch {
	case level < slog.LevelInfo: // Blue for DEBUG
		return "\x1b[34m" + "DEBUG" + "\x1b[0m"
	case level < slog.LevelWarn: // Green for INFO
		return "\x1b[32m" + "INFO " + "\x1b[0m"
	case level < slog.LevelError: // Yellow for WARN
		return "\x1b[33m" + "WARN " + "\x1b[0m"
	default: // Red for ERROR and higher
		return "\x1b[31m" + "ERROR" + "\x1b[0m"
	}
}

func (h *customConsoleHandler) Handle(_ context.Context, r slog.Record) error {
	_, err := fmt.Fprintln(os.Stdout, fmt.Errorf(
		"%v %v %v",
		r.Time.Format("2006/01/02 15:04:05"),
		getColouredLevelLabel(r.Level),
		r.Message,
	))
	return err
}
