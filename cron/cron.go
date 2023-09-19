package cron

import (
	"context"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
	"github.com/reugn/go-quartz/quartz"
	"github.com/rs/zerolog/log"
)

// RetentionPolicy
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
		r.app.Tracer(trace.SourceCron).Trace(ctx, "cron.RetentionPolicy", trace.WithError(err))
		log.Err(err).Msg("Failed to run app.RetentionPolicyRun")
	}
}

func (r RetentionPolicy) Key() int {
	return quartz.HashCode(r.Description())
}

// AttachmentOrphan
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

// Healthcheck
type Healthcheck struct {
	app core.App
	url string
}

func NewHealthcheck(app core.App, url string) Healthcheck {
	return Healthcheck{
		app: app,
		url: url,
	}
}

func (Healthcheck) Description() string {
	return "healthcheck"
}

func (r Healthcheck) Execute(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.url, nil)
	if err != nil {
		r.app.Tracer(trace.SourceCron).Trace(ctx, "cron.Healthcheck", trace.WithError(err))
		log.Err(err).Msg("Failed create HTTP request")
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		r.app.Tracer(trace.SourceCron).Trace(ctx, "cron.Healthcheck", trace.WithError(err))
		log.Err(err).Msg("Failed send HTTP request")
		return
	}
	res.Body.Close()
}

func (r Healthcheck) Key() int {
	return quartz.HashCode(r.Description())
}

// DatabaseVacuum
type DatabaseVacuum struct {
	app core.App
}

func NewDatabaseVacuum(app core.App) DatabaseVacuum {
	return DatabaseVacuum{
		app: app,
	}
}

func (DatabaseVacuum) Description() string {
	return "database-vacuum"
}

func (r DatabaseVacuum) Execute(ctx context.Context) {
	err := r.app.DatabaseVacuum(ctx)
	if err != nil {
		r.app.Tracer(trace.SourceCron).Trace(ctx, "cron.DatabaseVacuum", trace.WithError(err))
		log.Err(err).Msg("Failed to vacuum database")
		return
	}
}

func (r DatabaseVacuum) Key() int {
	return quartz.HashCode(r.Description())
}
