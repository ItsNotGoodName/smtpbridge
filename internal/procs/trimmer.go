package procs

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
	"github.com/ItsNotGoodName/smtpbridge/internal/events"
	"github.com/ItsNotGoodName/smtpbridge/internal/files"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/rs/zerolog/log"
)

func TrimStart(cc core.Context) error {
	ctx := cc.Context()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-events.PublishTrimStart(cc):
		if res {
			return nil
		}

		return fmt.Errorf("already trimming")
	}
}

func TrimmerBackground(ctx context.Context, app core.App, policy models.RetentionPolicy) {
	envDeletedC := make(chan core.EventEnvelopeDeleted, 1)
	envCreatedC := make(chan core.EventEnvelopeCreated, 1)
	evtTrimStart := make(chan core.EventTrimStart, 1)

	go trimmer(app.Context(ctx), policy, envCreatedC, envDeletedC, evtTrimStart)

	events.OnEnvelopeCreated(app, func(cc core.Context, evt core.EventEnvelopeCreated) {
		select {
		case <-envCreatedC:
		default:
		}

		select {
		case envCreatedC <- evt:
		default:
		}
	})

	events.OnEnvelopeDeleted(app, func(cc core.Context, evt core.EventEnvelopeDeleted) {
		select {
		case <-envDeletedC:
		default:
		}

		select {
		case envDeletedC <- evt:
		default:
		}
	})

	events.OnTrimStart(app, func(cc core.Context, evt core.EventTrimStart) {
		select {
		case evtTrimStart <- evt:
		default:
			select {
			case <-ctx.Done():
				return
			case evt.Response <- false:
			}
		}
	})
}

func trimmer(
	cc core.Context,
	policy models.RetentionPolicy,
	envCreatedC <-chan core.EventEnvelopeCreated,
	envDeletedC <-chan core.EventEnvelopeDeleted,
	evtTrimStart <-chan core.EventTrimStart,
) {
	ctx := cc.Context()
	ticker := time.NewTicker(30 * time.Minute)

	clean := func() {
		trimmerDeleteByAge(cc, policy)
		trimmerDeleteOrphanAttachments(cc)

		storage, err := StorageGet(cc)
		if err != nil {
			log.Err(err).Msg("Failed to get storage")
			return
		}

		trimmerDeleteByEnvelopeCount(cc, policy, storage)
		trimmerDeleteByAttachmentSize(cc, policy, storage)
	}
	clean()

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

			trimmerDeleteByEnvelopeCount(cc, policy, storage)
			trimmerDeleteByAttachmentSize(cc, policy, storage)
		case <-envDeletedC:
			trimmerDeleteOrphanAttachments(cc)
		case <-ticker.C:
			clean()
		case evt := <-evtTrimStart:
			clean()

			select {
			case <-ctx.Done():
				return
			case evt.Response <- true:
			}
		}
	}
}
func trimmerDeleteByAttachmentSize(cc core.Context, policy models.RetentionPolicy, storage models.Storage) {
	if policy.AttachmentSize == nil {
		return
	}
	attachmentSize := *policy.AttachmentSize

	if storage.AttachmentSize > attachmentSize {
		age := policy.AgeDate()
		log.Info().Time("age", age).Msg("Deleting attachment files by attachment size retention policy")
		err := files.DeleteFileUntilSize(cc, storage.AttachmentSize, attachmentSize, age)
		if err != nil {
			log.Err(err).Msg("Failed to delete attachment files by attachment size retention policy")
		}
	}
}

func trimmerDeleteByEnvelopeCount(cc core.Context, policy models.RetentionPolicy, storage models.Storage) {
	if policy.EnvelopeCount == nil {
		return
	}
	envelopeCount := *policy.EnvelopeCount

	if storage.EnvelopeCount > envelopeCount {
		age := policy.AgeDate()
		count, err := db.EnvelopeDeleteUntilCount(cc, envelopeCount, age)
		if err != nil {
			log.Err(err).Time("age", age).Int("keep", envelopeCount).Msg("Failed to envelopes by envelope count retention policy")
		} else {
			log.Info().Time("age", age).Int("keep", envelopeCount).Int64("deleted", count).Msg("Deleted envelopes by envelope count retention policy")
		}
	}
}

func trimmerDeleteByAge(cc core.Context, policy models.RetentionPolicy) {
	if policy.EnvelopeAge == nil {
		return
	}

	age := policy.AgeDate()
	count, err := db.EnvelopeDeleteOlderThan(cc, age)
	if err != nil {
		log.Err(err).Time("age", age).Msg("Failed to delete envelopes by age retention policy")
	} else {
		log.Info().Time("age", age).Int64("deleted", count).Msg("Deleted envelopes by age retention policy")
	}
}

func trimmerDeleteOrphanAttachments(cc core.Context) {
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
