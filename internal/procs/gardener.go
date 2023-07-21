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
	storageC := make(chan core.EventStorageRead, 1)
	envelopeC := make(chan core.EventEnvelopeDeleted, 1)

	go gardener(app.Context(ctx), policy, storageC, envelopeC)

	events.OnStorageRead(app, func(cc *core.Context, evt core.EventStorageRead) {
		select {
		case <-storageC:
		default:
		}

		select {
		case storageC <- evt:
		default:
		}
	})

	events.OnEnvelopeDeleted(app, func(cc *core.Context, evt core.EventEnvelopeDeleted) {
		select {
		case <-envelopeC:
		default:
		}

		select {
		case envelopeC <- evt:
		default:
		}
	})
}

func gardener(cc *core.Context, policy models.RetentionPolicy, storageC <-chan core.EventStorageRead, envelopeC <-chan core.EventEnvelopeDeleted) {
	ctx := cc.Context()
	ticker := time.NewTicker(60 * time.Minute)

	gardenerDeleteByAge(cc, policy)
	gardenerDeleteOrphanAttachments(cc)

	for {
		select {
		case <-ctx.Done():
			return
		case <-envelopeC:
			gardenerDeleteOrphanAttachments(cc)
		case evt := <-storageC:
			gardenerDeleteByStorage(cc, policy, evt.Storage)
		case <-ticker.C:
			storage, err := StorageGet(cc)
			if err != nil {
				log.Err(err).Msg("Failed to get storage")
				continue
			}
			gardenerDeleteByStorage(cc, policy, storage)
			gardenerDeleteByAge(cc, policy)
		}
	}
}
func gardenerDeleteByStorage(cc *core.Context, policy models.RetentionPolicy, storage models.Storage) {
	if policy.AttachmentSize != 0 && storage.AttachmentSize > policy.AttachmentSize {
		count := humanize.Bytes(uint64(storage.AttachmentSize - policy.AttachmentSize))
		log.Info().Str("count", count).Msg("Deleting attachment files by attachment size retention policy")

		err := files.DeleteFileUntilSize(cc, storage.AttachmentSize, policy.AttachmentSize)
		if err != nil {
			log.Err(err).Msg("Failed to delete attachment files by attachment size retention policy")
		}
	}

	if policy.EnvelopeCount != 0 && storage.EnvelopeCount > policy.EnvelopeCount {
		date := time.Now().Add(-policy.MinEnvelopeAge)
		count, err := db.EnvelopeDeleteUntilCount(cc, date, policy.EnvelopeCount)
		if err != nil {
			log.Err(err).Time("older-than", date).Int("keep", policy.EnvelopeCount).Msg("Failed to envelopes by envelope count retention policy")
		} else {
			log.Info().Time("older-than", date).Int("keep", policy.EnvelopeCount).Int64("deleted", count).Msg("Deleted envelopes by envelope count retention policy")
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
			log.Err(err).Time("older-than", date).Msg("Failed to delete envelopes by age retention policy")
		} else {
			log.Info().Time("older-than", date).Int64("deleted", count).Msg("Deleted envelopes by age retention policy")
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
