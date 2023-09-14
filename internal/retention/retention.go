package retention

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/repo"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
)

type FileStore interface {
	Remove(ctx context.Context, att models.Attachment) error
	Size(ctx context.Context) (int64, error)
	Trim(ctx context.Context, size int64, minAge time.Time) (int, error)
}

func DeleteAttachmentBySize(ctx context.Context, tracer trace.Tracer, fileStore FileStore, policy models.ConfigRetentionPolicy) (int, error) {
	if policy.AttachmentSize == nil {
		return 0, nil
	}
	maxAttachmentSize := *policy.AttachmentSize

	fileStoreSize, err := fileStore.Size(ctx)
	if err != nil {
		return 0, err
	}

	if fileStoreSize <= maxAttachmentSize {
		return 0, nil
	}

	age := policy.MinAgeTime()
	count, err := fileStore.Trim(ctx, maxAttachmentSize, age)
	if err != nil {
		return 0, err
	}

	tracer.Trace(ctx, "retention.attachment.size.delete", trace.WithKV("count", count))

	return count, nil
}

func DeleteEnvelopeByCount(ctx context.Context, tracer trace.Tracer, db database.Querier, policy models.ConfigRetentionPolicy) (int64, error) {
	if policy.EnvelopeCount == nil {
		return 0, nil
	}
	maxEnvelopeCount := *policy.EnvelopeCount

	repoEnvelopeCount, err := repo.EnvelopeCount(ctx, db)
	if err != nil {
		return 0, err
	}

	if repoEnvelopeCount <= maxEnvelopeCount {
		return 0, nil
	}

	age := policy.MinAgeTime()
	count, err := repo.EnvelopeTrim(ctx, db, age, maxEnvelopeCount)
	if err != nil {
		return 0, err
	}

	tracer.Trace(ctx, "retention.envelope.count.delete", trace.WithKV("count", count))

	return count, nil
}

func DeleteEnvelopeByAge(ctx context.Context, tracer trace.Tracer, db database.Querier, policy models.ConfigRetentionPolicy) (int64, error) {
	if policy.EnvelopeAge == nil {
		return 0, nil
	}

	age := policy.EnvelopeAgeTime()
	count, err := repo.EnvelopeTrim(ctx, db, age, 0)
	if err != nil {
		return 0, err
	}

	tracer.Trace(ctx, "retention.envelope.age.delete", trace.WithKV("count", count))

	return count, nil
}

func DeleteOrphanAttachments(ctx context.Context, tracer trace.Tracer, db database.Querier, fileStore FileStore) error {
	var lastID int64 = -1
	for {
		atts, err := repo.AttachmentListOrphan(ctx, db, 10)
		if err != nil {
			return err
		}
		if len(atts) == 0 {
			return nil
		}

		for _, a := range atts {
			// sanity check
			if a.ID == lastID {
				return fmt.Errorf("infinite loop detected")
			}
			lastID = a.ID

			if err := fileStore.Remove(ctx, a); err != nil {
				return err
			}

			if err := repo.AttachmentRemove(ctx, db, a.ID); err != nil {
				return err
			}
		}

		tracer.Trace(ctx, "retention.attachment.orphan.delete", trace.WithKV("count", len(atts)))
	}
}

func DeleteTraceByAge(ctx context.Context, tracer trace.Tracer, db database.Querier, policy models.ConfigRetentionPolicy) (int64, error) {
	if policy.TraceAge == nil {
		return 0, nil
	}

	age := time.Now().Add(-*policy.TraceAge)
	count, err := repo.TraceTrim(ctx, db, age)
	if err != nil {
		return 0, err
	}

	tracer.Trace(ctx, "retention.trace.age.delete", trace.WithKV("count", count))

	return count, nil
}
