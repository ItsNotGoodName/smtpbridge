package core

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
)

type Bus interface {
	EnvelopeCreated(ctx context.Context, id int64)
	OnEnvelopeCreated(func(ctx context.Context, evt models.EventEnvelopeCreated) error) func()
	EnvelopeDeleted(ctx context.Context)
	OnEnvelopeDeleted(func(ctx context.Context, evt models.EventEnvelopeDeleted) error) func()
	MailmanEnqueued(ctx context.Context)
	OnMailmanEnqueued(func(ctx context.Context, evt models.EventMailmanEnqueued) error) func()
}

type App interface {
	AttachmentGet(ctx context.Context, id int64) (models.Attachment, error)
	AttachmentList(ctx context.Context, page pagination.Page, filter models.DTOAttachmentListRequest) (models.DTOAttachmentListResult, error)
	AttachmentOrphanDelete(ctx context.Context, tracer trace.Tracer) error
	AuthHTTPAnonymous() bool
	AuthHTTPLogin(ctx context.Context, username, password string) (models.User, error)
	AuthSMTPLogin(ctx context.Context, username, password string) error
	EndpointList(ctx context.Context) ([]models.Endpoint, error)
	EndpointTest(ctx context.Context, id int64) error
	EnvelopeCreate(ctx context.Context, msg models.DTOMessageCreate, datts []models.DTOAttachmentCreate) (int64, error)
	EnvelopeDelete(ctx context.Context, id int64) error
	EnvelopeDrop(ctx context.Context) error
	EnvelopeGet(ctx context.Context, id int64) (models.Envelope, error)
	EnvelopeList(ctx context.Context, page pagination.Page, req models.DTOEnvelopeListRequest) (models.DTOEnvelopeListResult, error)
	EnvelopeSend(ctx context.Context, envelopeID int64, endpointID int64) error
	EnvelopeCount(ctx context.Context) (int, error)
	MessageHTMLGet(ctx context.Context, id int64) (string, error)
	RetentionPolicyGet(ctx context.Context) models.ConfigRetentionPolicy
	RetentionPolicyRun(ctx context.Context, trace trace.Tracer) error
	RuleCreate(ctx context.Context, req models.DTORuleCreate) (int64, error)
	RuleDelete(ctx context.Context, id int64) error
	RuleEndpointsList(ctx context.Context) ([]models.RuleEndpoints, error)
	RuleExpressionCheck(ctx context.Context, expression string) error
	RuleGet(ctx context.Context, id int64) (models.Rule, error)
	RuleEndpointsGet(ctx context.Context, id int64) (models.RuleEndpoints, error)
	RuleList(ctx context.Context) ([]models.Rule, error)
	RuleUpdate(ctx context.Context, req models.DTORuleUpdate) error
	StorageGet(ctx context.Context) (models.Storage, error)
	TraceDrop(ctx context.Context) error
	TraceList(ctx context.Context, page pagination.Page, req models.DTOTraceListRequest) (models.DTOTraceListResult, error)
	Tracer(source string) trace.Tracer
	MailmanEnqueue(ctx context.Context, envelopeID int64) error
	MailmanDequeue(ctx context.Context) (*models.Envelope, error)
}
