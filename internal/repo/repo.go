package repo

import (
	"database/sql"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/go-jet/jet/v2/qrm"
)

var ErrNoRows = qrm.ErrNoRows

func one(res sql.Result) error {
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return qrm.ErrNoRows
	}

	return nil
}

func muhTypeAffinity(date models.Time) string {
	return "\"" + date.Time().Format(time.RFC3339) + "\""
}
