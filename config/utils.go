package config

import (
	"os"

	"github.com/alecthomas/kong"
)

func resolve(paths []string) (string, error) {
	for _, path := range paths {
		f, err := os.Open(kong.ExpandPath(path))
		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) {
				continue
			}

			return "", err
		}
		f.Close()

		return path, nil
	}

	return "", nil
}
