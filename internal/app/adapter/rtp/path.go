package rtp

import (
	"os"
	"os/user"
	"path/filepath"
)

var ownDir string

func OwnDir() string {
	if ownDir != "" {
		return ownDir
	}

	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	ownDir = filepath.Join(u.HomeDir, ".comhel")
	err = os.MkdirAll(ownDir, 0700)
	if err != nil {
		panic(err)
	}

	return ownDir
}
