package gke

import (
	"log/slog"
	"os"
)

type LogLevel slog.Level

const (
	LogDebug   = LogLevel(slog.LevelDebug)
	LogInfo    = LogLevel(slog.LevelInfo)
	LogWarning = LogLevel(slog.LevelWarn)
	LogError   = LogLevel(slog.LevelError)
)

var logLevel = new(slog.LevelVar)
var log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	Level: logLevel, // LevelDebug, LevelInfo, LevelWarn, LevelError
}))

func setLogLevel(uroven LogLevel) {
	logLevel.Set(slog.Level(uroven))
}
