package logattr

import (
	"fmt"
	"log/slog"
)

func Any(name string, val any) slog.Attr {
	return slog.String(name, fmt.Sprintf("%#v", val))
}

func Caller(s string) slog.Attr {
	return slog.String("caller", s)
}
