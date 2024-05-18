package logger

import (
	"log/slog"
	"os"
)

var levelInfo map[string]slog.Level = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"":      slog.LevelInfo,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

var log *slog.Logger

func init() {
	log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     levelInfo[os.Getenv("LOGGING_LEVEL")],
	}))
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	Error(msg, args...)
	os.Exit(1)
}
