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

var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     levelInfo[os.Getenv("LOGGING_LEVEL")],
	}))
}
