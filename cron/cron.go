package cron

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
	"github.com/reugn/go-quartz/quartz"
	"github.com/rs/zerolog/log"
)

type RetentionPolicy struct {
	app core.App
}

func NewRetentionPolicy(app core.App) RetentionPolicy {
	return RetentionPolicy{
		app: app,
	}
}

func (RetentionPolicy) Description() string {
	return "retention-policy"
}

func (r RetentionPolicy) Execute(ctx context.Context) {
	err := r.app.RetentionPolicyRun(ctx, r.app.Tracer(trace.SourceCron))
	if err != nil {
		log.Err(err).Msg("Failed to run app.RetentionPolicyRun")
	}
}

func (r RetentionPolicy) Key() int {
	return quartz.HashCode(r.Description())
}

type AttachmentOrphan struct {
	app core.App
}

func NewAttachmentOrphan(app core.App) AttachmentOrphan {
	return AttachmentOrphan{
		app: app,
	}
}

func (AttachmentOrphan) Description() string {
	return "attachment-orphan"
}

func (r AttachmentOrphan) Execute(ctx context.Context) {
	err := r.app.AttachmentOrphanDelete(ctx, r.app.Tracer(trace.SourceCron))
	if err != nil {
		log.Err(err).Msg("Failed to run app.AttachmentTrim")
	}
}

func (r AttachmentOrphan) Key() int {
	return quartz.HashCode(r.Description())
}
