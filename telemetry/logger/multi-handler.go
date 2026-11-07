package logger

import (
	"context"
	"log/slog"
)

type multiHandler struct {
	handlers []slog.Handler
}

// newMultiHandler creates a new MultiHandler instance
func newMultiHandler(
	handlers ...slog.Handler,
) *multiHandler {
	return &multiHandler{
		handlers: handlers,
	}
}

// Enabled returns if this handler is enabled.
// If any of the underlying handlers are enabled, the multiHandler is enabled.
func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle Call Handle on all underlying handlers
func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		// Only call Handle if the individual handler is enabled for this record's level
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}
