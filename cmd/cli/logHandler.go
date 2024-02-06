package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/fatih/color"
)

type LogHandler struct {
	handler slog.Handler
	logger  *log.Logger
}

func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

func (h *LogHandler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

func (h *LogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()
	switch r.Level {
	case slog.LevelDebug:
		level = color.GreenString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	attrStr := ""
	if r.NumAttrs() > 0 {
		fields := make(map[string]interface{}, r.NumAttrs())
		r.Attrs(func(a slog.Attr) bool {
			fields[a.Key] = a.Value.Any()

			return true
		})

		b, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
		attrStr = string(b)
	}

	timeStr := r.Time.Format("[15:04:05.000]")

	h.logger.Println(timeStr, level, color.WhiteString(r.Message), color.WhiteString(attrStr))

	return nil
}

func NewLogHandler(
	out io.Writer,
	logLevel string,
) *LogHandler {
	opts := slog.HandlerOptions{Level: parseLogLevel(logLevel)}
	h := &LogHandler{
		handler: slog.NewTextHandler(out, &opts),
		logger:  log.New(out, "", 0),
	}

	return h
}

func parseLogLevel(logLevel string) slog.Level {
	switch strings.ToLower(logLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		fmt.Printf("Invalid log level value '%s'. Using the default log level INFO.\n", logLevel)
		return slog.LevelInfo
	}
}
