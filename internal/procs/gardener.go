package procs

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
	"github.com/ItsNotGoodName/smtpbridge/internal/events"
	"github.com/ItsNotGoodName/smtpbridge/internal/files"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
)

func GardenerStart(ctx context.Context, app core.App, policy models.RetentionPolicy) {
	envDeletedC := make(chan core.EventEnvelopeDeleted, 1)
	envCreatedC := make(chan core.EventEnvelopeCreated, 1)

	go gardener(app.Context(ctx), policy, envCreatedC, envDeletedC)

	events.OnEnvelopeCreated(app, func(cc *core.Context, evt core.EventEnvelopeCreated) {
		select {
		case <-envCreatedC:
		default:
		}

		select {
		case envCreatedC <- evt:
		default:
		}
	})

	events.OnEnvelopeDeleted(app, func(cc *core.Context, evt core.EventEnvelopeDeleted) {
		select {
		case <-envDeletedC:
		default:
		}

		select {
		case envDeletedC <- evt:
		default:
		}
	})
}

func gardener(cc *core.Context, policy models.RetentionPolicy, envCreatedC <-chan core.EventEnvelopeCreated, envDeletedC <-chan core.EventEnvelopeDeleted) {
	ctx := cc.Context()
	ticker := time.NewTicker(30 * time.Minute)

	gardenerDeleteByAge(cc, policy)

	gardenerDeleteOrphanAttachments(cc)

	for {
		select {
		case <-ctx.Done():
			return
		case <-envCreatedC:
			storage, err := StorageGet(cc)
			if err != nil {
				log.Err(err).Msg("Failed to get storage")
				continue
			}

			gardenerDeleteByEnvelopeCount(cc, policy, storage)
			gardenerDeleteByAttachmentSize(cc, policy, storage)
		case <-envDeletedC:
			gardenerDeleteOrphanAttachments(cc)
		case <-ticker.C:
			gardenerDeleteByAge(cc, policy)
		}
	}
}
func gardenerDeleteByAttachmentSize(cc *core.Context, policy models.RetentionPolicy, storage models.Storage) {
	if policy.AttachmentSize != 0 && storage.AttachmentSize > policy.AttachmentSize {
		count := humanize.Bytes(uint64(storage.AttachmentSize - policy.AttachmentSize))
		log.Info().Str("count", count).Msg("Deleting attachment files by attachment size retention policy")

		err := files.DeleteFileUntilSize(cc, storage.AttachmentSize, policy.AttachmentSize)
		if err != nil {
			log.Err(err).Msg("Failed to delete attachment files by attachment size retention policy")
		}
	}
}

func gardenerDeleteByEnvelopeCount(cc *core.Context, policy models.RetentionPolicy, storage models.Storage) {
	if policy.EnvelopeCount != 0 && storage.EnvelopeCount > policy.EnvelopeCount {
		date := time.Now().Add(-policy.MinEnvelopeAge)
		count, err := db.EnvelopeDeleteUntilCount(cc, policy.EnvelopeCount, date)
		if err != nil {
			log.Err(err).Time("age", date).Int("keep", policy.EnvelopeCount).Msg("Failed to envelopes by envelope count retention policy")
		} else {
			log.Info().Time("age", date).Int("keep", policy.EnvelopeCount).Int64("deleted", count).Msg("Deleted envelopes by envelope count retention policy")
		}
	}
}

func gardenerDeleteByAge(cc *core.Context, policy models.RetentionPolicy) {
	if policy.EnvelopeAge != 0 {
		date := time.Now().Add(-policy.EnvelopeAge)
		if policy.MinEnvelopeAge > policy.EnvelopeAge {
			date.Add(-policy.MinEnvelopeAge)
		}
		count, err := db.EnvelopeDeleteOlderThan(cc, date)
		if err != nil {
			log.Err(err).Time("age", date).Msg("Failed to delete envelopes by age retention policy")
		} else {
			log.Info().Time("age", date).Int64("deleted", count).Msg("Deleted envelopes by age retention policy")
		}
	}
}

func gardenerDeleteOrphanAttachments(cc *core.Context) {
	for {
		atts, err := db.EnvelopeAttachmentListOrphan(cc, 10)
		if err != nil {
			log.Err(err).Msg("Failed to list orphan attachments")
			return
		}
		if len(atts) == 0 {
			return
		}

		for _, a := range atts {
			log.Info().Int64("id", a.ID).Msg("Deleting orphan attachment")
			err := db.EnvelopeAttachmentDelete(cc, a)
			if err != nil {
				log.Err(err).Msg("Failed to delete orphan attachment")
				return
			}
		}
	}
}
