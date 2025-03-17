package env

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/FaisalBudiono/comhel/internal/app/adapter/rtp"
	"github.com/joho/godotenv"
)

type spec struct {
	DevMode  bool
	LogLevel string
}

var s spec

func Get() spec {
	return s
}

func Bind() error {
	err := godotenv.Load(envFilePaths()...)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	setENV()

	return nil
}

func setENV() {
	s = spec{
		DevMode:  strings.ToLower(os.Getenv("COMHEL_DEV_MODE")) == "true",
		LogLevel: strings.ToLower(os.Getenv("COMHEL_LOG_LEVEL")),
	}
}

func envFilePaths() []string {
	envName := ".env"

	return []string{
		envName,
		filepath.Join(rtp.OwnDir(), envName),
	}
}
