package procs

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/rs/zerolog/log"
)

func VacuumerBackground(ctx context.Context, app core.App) {
	go vacuum(app.SystemContext(ctx))
}

func vacuum(cc core.Context) {
	ctx := cc.Context()
	ticker := time.NewTicker(24 * time.Hour)

	run := func() {
		log.Info().Msg("Vacuuming database")
		err := DatabaseVacuum(cc)
		if err != nil {
			log.Err(err).Msg("Failed to vacuum")
		}
	}
	run()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}
