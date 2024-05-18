package logger

import (
	"log/slog"
	"os"
	"sync"
)

var once sync.Once
var log *slog.Logger

func init() {
	once.Do(func() {
		log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}))
	})
}

func Logger() *slog.Logger {
	return log
}
