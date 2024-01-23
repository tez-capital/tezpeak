package util

import (
	"log/slog"
	"os"
)

// / ideally move somewhere else where it makes sense
func InitLog(logLevel string) {
	l := slog.LevelInfo
	switch logLevel {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		slog.Error("invalid log level")
		os.Exit(1)
	}

	// replace next 4 lines with SetLogLoggerLevel in go 1.22
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: l,
	})
	slog.SetDefault(slog.New(textHandler))
}
