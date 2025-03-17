package log

import (
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
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
	f, err := os.OpenFile(logPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	slogHandler := slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level: logLevel(),
	})

	return slog.New(slogHandler), nil
}

func logPath() string {
	devMode := strings.ToLower(os.Getenv("DEV_MODE")) == "true"
	if devMode {
		return "./logs/logs.log"
	}

	return prodLogPath()
}

func prodLogPath() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	dir := filepath.Join(u.HomeDir, ".comhel")
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		panic(err)
	}

	return filepath.Join(dir, "logs.log")
}

func logLevel() slog.Leveler {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "debug" {
		return slog.LevelDebug
	}

	return slog.LevelWarn
}
