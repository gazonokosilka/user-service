package slog

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/mgutz/ansi"
)

// CustomHandler custom handler for beauty terminal.
type CustomHandler struct {
	slog.Handler
	attrs []slog.Attr
}

// NewCustomHandler create new CustomHandler.
func NewCustomHandler(handler slog.Handler) *CustomHandler {
	return &CustomHandler{Handler: handler}
}

// Handle processing log record.
func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()
	var levelColor string

	switch level {
	case "INFO":
		levelColor = ansi.Green
	case "WARN":
		levelColor = ansi.Yellow
	case "ERROR":
		levelColor = ansi.Red
	case "DEBUG":
		levelColor = ansi.Blue
	default:
		levelColor = ansi.White
	}

	resetColor := ansi.Reset

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("%s%s%s %s: %s",
		levelColor, level, resetColor,
		r.Time.Format(time.RFC3339),
		r.Message))

	r.Attrs(func(a slog.Attr) bool {
		key := a.Key
		val := a.Value.String()

		if key == "op" {

			msg.WriteString(fmt.Sprintf(" %s%s=%s%s", ansi.Magenta, key, val, ansi.Reset))
		} else {
			msg.WriteString(fmt.Sprintf(" %s=%s", key, val))
		}
		return true
	})

	for _, attr := range h.attrs {
		msg.WriteString(fmt.Sprintf(" %s=%v", attr.Key, attr.Value.Any()))
	}

	for _, attr := range h.attrs {
		key := attr.Key
		val, ok := attr.Value.Any().(string)
		if !ok {
			val = fmt.Sprint(attr.Value.Any())
		}

		if key == "op" {
			msg.WriteString(fmt.Sprintf(" %s%s=%s%s", ansi.Magenta, key, val, ansi.Reset))
		} else {
			msg.WriteString(fmt.Sprintf(" %s=%s", key, val))
		}
	}

	fmt.Println()
	fmt.Println(msg.String())

	return nil
}

// WithAttrs new handler with additional attributes
func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &CustomHandler{Handler: h.Handler, attrs: newAttrs}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return h
}
