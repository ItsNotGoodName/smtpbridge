package app

import (
	"errors"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/repo"
)

func checkErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repo.ErrNoRows) {
		return models.ErrNotFound
	}

	return err
}
