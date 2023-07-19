package procs

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
	"github.com/ItsNotGoodName/smtpbridge/internal/events"
	"github.com/ItsNotGoodName/smtpbridge/internal/files"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
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
		err := files.DeleteFileUntilSize(cc, storage.AttachmentSize, policy.AttachmentSize)
		if err != nil {
			log.Err(err).Msg("Failed to trim files")
		}
	}

	if policy.EnvelopeCount != 0 && storage.EnvelopeCount > policy.EnvelopeCount {
		_, err := db.EnvelopeDeleteUntilCount(cc, policy.EnvelopeCount)
		if err != nil {
			log.Err(err).Msg("Failed to trim envelopes")
		}
	}
}

func gardenerDeleteByAge(cc *core.Context, policy models.RetentionPolicy) {
	if policy.EnvelopeAge != 0 {
		db.EnvelopeDeleteOlderThan(cc, time.Now().Add(-policy.EnvelopeAge))
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
			err := db.EnvelopeAttachmentDelete(cc, a)
			if err != nil {
				log.Err(err).Msg("Failed to delete orphan attachment")
				return
			}
		}
	}
}
