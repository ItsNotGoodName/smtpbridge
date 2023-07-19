package procs

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
)

func DatabaseVacuum(cc *core.Context) error {
	return db.Vacuum(cc)
}
