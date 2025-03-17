package jsonconfig

import "os"

var filename = "./.comhelconfig.json"

func openFile() (*os.File, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}
