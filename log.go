package gke

import (
	"log/slog"
	"os"
)

type LogLevel slog.Level

const (
	LogDebug   = slog.LevelDebug
	LogInfo    = slog.LevelInfo
	LogWarning = slog.LevelWarn
	LogError   = slog.LevelError
)

var logLevel = new(slog.LevelVar)
var log = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	Level: logLevel, // LevelDebug, LevelInfo, LevelWarn, LevelError
})

func setLogLevel(uroven LogLevel) {
	logLevel.Set(slog.Level(uroven))
}
