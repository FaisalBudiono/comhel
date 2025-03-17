package log

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/FaisalBudiono/comhel/internal/app/adapter/env"
	"github.com/FaisalBudiono/comhel/internal/app/adapter/rtp"
)

var logger *slog.Logger

func Logger() *slog.Logger {
	return logger
}

func SetDefault(l *slog.Logger) {
	logger = l
}

func NewLogger() (*slog.Logger, error) {
	f, err := os.OpenFile(logPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	slogHandler := slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level: logLevel(),
	})

	return slog.New(slogHandler), nil
}

var logFilename = "logs.log"

func logPath() string {
	if env.Get().DevMode {
		return filepath.Join("./logs", logFilename)
	}

	return prodLogPath()
}

func prodLogPath() string {
	return filepath.Join(rtp.OwnDir(), logFilename)
}

func logLevel() slog.Leveler {
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch level {
	case "debug":
		return slog.LevelDebug
	default:
		return slog.LevelWarn
	}
}
