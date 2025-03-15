package env

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func Bind() error {
	err := godotenv.Load()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}
