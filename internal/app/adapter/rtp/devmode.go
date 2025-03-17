package rtp

import (
	"os"
	"strings"
)

func DevMode() bool {
	return strings.ToLower(os.Getenv("DEV_MODE")) == "true"
}
