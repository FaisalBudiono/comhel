package env

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/FaisalBudiono/comhel/internal/app/adapter/rtp"
	"github.com/joho/godotenv"
)

func Bind() error {
	err := godotenv.Load(envFilePath())
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func envFilePath() string {
	envName := ".env"

	if rtp.DevMode() {
		return envName
	}

	return filepath.Join(rtp.OwnDir(), envName)
}
