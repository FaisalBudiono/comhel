package logattr

import "log/slog"

func Caller(s string) slog.Attr {
	return slog.String("caller", s)
}
