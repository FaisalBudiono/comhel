package log

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/FaisalBudiono/comhel/internal/app/adapter/rtp"
	"github.com/FaisalBudiono/comhel/internal/app/core/util/env"
)

func New() (*slog.Logger, error) {
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
	level := env.Get().LogLevel

	switch level {
	case "debug":
		return slog.LevelDebug
	default:
		return slog.LevelWarn
	}
}
