package log

import (
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func Logger() *slog.Logger {
	return logger
}

func SetDefault(l *slog.Logger) {
	logger = l
}

func NewLogger() (*slog.Logger, error) {
	f, err := os.OpenFile("./logs/comhel.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	slogHandler := slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level: logLevel(),
	})

	return slog.New(slogHandler), nil
}

func logLevel() slog.Leveler {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "debug" {
		return slog.LevelDebug
	}

	return slog.LevelWarn
}
