package logger

import (
	"log/slog"
	"os"
)

func New(level string) *slog.Logger {
	var programLevel = new(slog.LevelVar) // Info by default
	if err := programLevel.UnmarshalText([]byte(level)); err != nil {
		panic("invalid log level")
	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
	logger := slog.Default()
	logger.Info("logger initialized", "level", level)
	return logger
}
