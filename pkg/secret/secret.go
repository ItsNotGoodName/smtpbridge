package secret

import (
	"crypto/rand"
	"errors"
	"os"
)

func GetOrCreate(filePath string) ([]byte, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		b = make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return nil, err
		}

		os.WriteFile(filePath, b, 0600)
	}

	return b, nil
}
