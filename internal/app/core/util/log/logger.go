package log

import "log/slog"

var logger *slog.Logger

func Logger() *slog.Logger {
	return logger
}

func SetDefault(l *slog.Logger) {
	logger = l
}
