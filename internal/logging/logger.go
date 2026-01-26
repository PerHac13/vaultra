package logging

import (
	"io"
	"log/slog"
	"os"
)


type Logger struct {
	*slog.Logger
}

func NewLogger(level slog.Level, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(output, opts)

	return &Logger{
		Logger: slog.New(handler),
	}
}

func NewDefaultLogger() *Logger {
	return NewLogger(slog.LevelInfo, os.Stdout)
}