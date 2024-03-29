package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/auth"
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	"github.com/ItsNotGoodName/smtpbridge/internal/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/repo"
	"github.com/ItsNotGoodName/smtpbridge/internal/retention"
	"github.com/ItsNotGoodName/smtpbridge/internal/rule"
	"github.com/ItsNotGoodName/smtpbridge/internal/trace"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
)

type FileStore interface {
	Create(ctx context.Context, att models.Attachment, data io.Reader) error
	Remove(ctx context.Context, att models.Attachment) error
	Size(ctx context.Context) (int64, error)
	Trim(ctx context.Context, size int64, minAge time.Time) (int, error)
	Reader(ctx context.Context, att models.Attachment) (io.ReadCloser, error)
	Path(ctx context.Context, att models.Attachment) (string, error)
}

type App struct {
	db               database.Querier
	fileStore        FileStore
	bus              core.Bus
	config           *models.Config
	endpointFactory  endpoint.Factory
	webTestFileStore WebTestFileStore
}

func New(
	db database.Querier,
	fileStore FileStore,
	bus core.Bus,
	config *models.Config,
	endpointFactory endpoint.Factory,
	webTestFileStore WebTestFileStore,
) (App, func()) {
	a := App{
		db:               db,
		fileStore:        fileStore,
		bus:              bus,
		config:           config,
		endpointFactory:  endpointFactory,
		webTestFileStore: webTestFileStore,
	}

	return a, a.init()
}

func (a App) init() func() {
	return closers(
		a.bus.OnEnvelopeCreated(func(ctx context.Context, evt models.EventEnvelopeCreated) error {
			return a.MailmanEnqueue(ctx, evt.ID)
		}),
		a.bus.OnEnvelopeDeleted(func(ctx context.Context, evt models.EventEnvelopeDeleted) error {
			return a.AttachmentOrphanDelete(ctx, a.Tracer(trace.SourceApp))
		}),
	)
}

func (a App) DatabaseVacuum(ctx context.Context) error {
	return repo.Vacuum(ctx, a.db)
}

// - Auth

func (a App) AuthSMTPAnonymous() bool {
	return a.config.AuthSMTP.Anonymous
}

// AuthHTTPAnonymous checks if anonymous access is allowed.
func (a App) AuthHTTPAnonymous() bool {
	return a.config.AuthHTTP.Anonymous
}

func (a App) AuthHTTPLogin(ctx context.Context, username, password string) (models.User, error) {
	if a.config.AuthHTTP.Anonymous || auth.Check(a.config.AuthHTTP, username, password) {
		return models.User{}, nil
	}

	return models.User{}, models.ErrAuthInvalid
}

func (a App) AuthSMTPLogin(ctx context.Context, username, password string) error {
	if a.config.AuthSMTP.Anonymous || auth.Check(a.config.AuthSMTP, username, password) {
		return nil
	}

	return models.ErrAuthInvalid
}

// - Endpoint

func (a App) EndpointCreate(ctx context.Context, req models.DTOEndpointCreate) (int64, error) {
	end, err := endpoint.New(a.endpointFactory, req)
	if err != nil {
		return 0, err
	}

	id, err := repo.EndpointCreate(ctx, a.db, end)
	return id, repoErr(err)
}

func (a App) EndpointGet(ctx context.Context, id int64) (models.Endpoint, error) {
	end, err := repo.EndpointGet(ctx, a.db, id)
	return end, repoErr(err)
}

func (a App) EndpointList(ctx context.Context) ([]models.Endpoint, error) {
	return repo.EndpointList(ctx, a.db)
}

func (a App) EndpointUpdate(ctx context.Context, req models.DTOEndpointUpdate) error {
	end, err := repo.EndpointGet(ctx, a.db, req.ID)
	if err != nil {
		return repoErr(err)
	}

	end, err = endpoint.Update(a.endpointFactory, end, req)
	if err != nil {
		return err
	}

	err = repo.EndpointUpdate(ctx, a.db, end)
	return repoErr(err)
}

func (a App) EndpointDelete(ctx context.Context, id int64) error {
	end, err := repo.EndpointGet(ctx, a.db, id)
	if err != nil {
		return repoErr(err)
	}

	err = endpoint.Delete(end)
	if err != nil {
		return err
	}

	err = repo.EndpointDelete(ctx, a.db, id)
	return repoErr(err)
}

func (a App) EndpointTest(ctx context.Context, id int64) error {
	e, err := repo.EndpointGet(ctx, a.db, id)
	if err != nil {
		return repoErr(err)
	}

	end, err := a.endpointFactory.Build(e)
	if err != nil {
		return err
	}

	file, err := a.webTestFileStore.File()
	if err != nil {
		return err
	}
	defer file.Close()

	datt, err := envelope.NewDataAttachment("Test Attachment", file)
	if err != nil {
		return err
	}
	file.Close()

	env := envelope.New(envelope.NewMessage(models.DTOMessageCreate{
		Subject: "Test Subject",
		Text:    "Test Body",
		From:    "from-test@example.com",
		To:      []string{"to-test@example.com"},
		Date:    time.Now(),
	}), datt.Attachment)

	return end.Send(ctx, a.webTestFileStore, env)
}

// - Envelope

func (a App) EnvelopeCreate(ctx context.Context, dtoMsg models.DTOMessageCreate, dtoDatts []models.DTOAttachmentCreate) (int64, error) {
	msg := envelope.NewMessage(dtoMsg)

	atts := make([]models.Attachment, 0, len(dtoDatts))
	datts := make([]models.DataAttachment, 0, len(dtoDatts))
	for _, dc := range dtoDatts {
		datt, err := envelope.NewDataAttachment(dc.Name, dc.Data)
		if err != nil {
			return 0, err
		}

		atts = append(atts, datt.Attachment)
		datts = append(datts, datt)
	}

	id, err := repo.EnvelopeCreate(ctx, a.db, msg, atts)
	if err != nil {
		return 0, err
	}

	atts, err = repo.AttachmentListByMessage(ctx, a.db, id)
	if err != nil {
		return 0, err
	}
	if len(atts) < len(datts) {
		return 0, fmt.Errorf("invalid number of attachments from repo")
	}

	for i := range atts {
		err := a.fileStore.Create(ctx, atts[i], datts[i])
		if err != nil {
			return 0, err
		}
	}

	a.bus.EnvelopeCreated(ctx, id)

	return id, nil
}

func (a App) EnvelopeCount(ctx context.Context) (int, error) {
	return repo.EnvelopeCount(ctx, a.db)
}

func (a App) EnvelopeGet(ctx context.Context, id int64) (models.Envelope, error) {
	res, err := repo.EnvelopeGet(ctx, a.db, id)
	return res, repoErr(err)
}

func (a App) EnvelopeList(ctx context.Context, page pagination.Page, req models.DTOEnvelopeListRequest) (models.DTOEnvelopeListResult, error) {
	return repo.EnvelopeList(ctx, a.db, page, req)
}

func (a App) EnvelopeDelete(ctx context.Context, id int64) error {
	if err := repo.EnvelopeDelete(ctx, a.db, id); err != nil {
		return repoErr(err)
	}

	a.bus.EnvelopeDeleted(ctx)

	return nil
}

func (a App) EnvelopeDrop(ctx context.Context) error {
	count, err := repo.EnvelopeDrop(ctx, a.db)
	if err != nil {
		return err
	}

	if count > 0 {
		a.bus.EnvelopeDeleted(ctx)
	}

	return nil
}

func (a App) EnvelopeSend(ctx context.Context, envelopeID int64, endpointID int64) error {
	env, err := repo.EnvelopeGet(ctx, a.db, envelopeID)
	if err != nil {
		return repoErr(err)
	}

	endModel, err := repo.EndpointGet(ctx, a.db, endpointID)
	if err != nil {
		return repoErr(err)
	}

	end, err := a.endpointFactory.Build(endModel)
	if err != nil {
		return err
	}

	return end.Send(ctx, a.fileStore, env)
}

// - Rule

func (a App) RuleCreate(ctx context.Context, req models.DTORuleCreate) (int64, error) {
	r, err := rule.New(req)
	if err != nil {
		return 0, err
	}

	return repo.RuleCreate(ctx, a.db, r, req.Endpoints)
}

func (a App) RuleUpdate(ctx context.Context, req models.DTORuleUpdate) error {
	r, err := repo.RuleGet(ctx, a.db, req.ID)
	if err != nil {
		return repoErr(err)
	}

	r, err = rule.Update(r, req)
	if err != nil {
		return err
	}

	err = repo.RuleUpdate(ctx, a.db, r)
	if err != nil {
		return repoErr(err)
	}

	if req.Endpoints != nil {
		err := repo.RuleEndpointsSet(ctx, a.db, r.ID, *req.Endpoints)
		if err != nil {
			return repoErr(err)
		}
	}

	return nil
}

func (a App) RuleGet(ctx context.Context, id int64) (models.Rule, error) {
	res, err := repo.RuleGet(ctx, a.db, id)
	return res, repoErr(err)
}

func (a App) RuleList(ctx context.Context) ([]models.Rule, error) {
	return repo.RuleList(ctx, a.db)
}

func (a App) RuleDelete(ctx context.Context, id int64) error {
	r, err := repo.RuleGet(ctx, a.db, id)
	if err != nil {
		return repoErr(err)
	}

	err = rule.Delete(r)
	if err != nil {
		return err
	}

	err = repo.RuleDelete(ctx, a.db, id)
	return repoErr(err)
}

func (App) RuleExpressionCheck(ctx context.Context, expression string) error {
	tmpl, err := rule.TemplateBuild(expression)
	if err != nil {
		return err
	}

	_, err = rule.TemplateRun(tmpl, models.Envelope{Message: models.Message{}, Attachments: []models.Attachment{}})
	if err != nil {
		return err
	}

	return nil
}

func (a App) RuleEndpointsGet(ctx context.Context, id int64) (models.RuleEndpoints, error) {
	res, err := repo.RuleEndpointsGet(ctx, a.db, id)
	return res, repoErr(err)
}

func (a App) RuleEndpointsList(ctx context.Context) ([]models.RuleEndpoints, error) {
	return repo.RuleEndpointsList(ctx, a.db)
}

// - Attachment

func (a App) AttachmentGet(ctx context.Context, id int64) (models.Attachment, error) {
	res, err := repo.AttachmentGet(ctx, a.db, id)
	return res, repoErr(err)
}

func (a App) AttachmentList(ctx context.Context, page pagination.Page, req models.DTOAttachmentListRequest) (models.DTOAttachmentListResult, error) {
	return repo.AttachmentList(ctx, a.db, page, req)
}

func (a App) AttachmentOrphanDelete(ctx context.Context, tracer trace.Tracer) error {
	return retention.DeleteOrphanAttachments(ctx, tracer, a.db, a.fileStore)
}

// - Message

func (a App) MessageHTMLGet(ctx context.Context, id int64) (string, error) {
	res, err := repo.MessageHTMLGet(ctx, a.db, id)
	return res, repoErr(err)
}

// - Storage

func (a App) StorageGet(ctx context.Context) (models.Storage, error) {
	attachmentCount, err := repo.AttachmentCount(ctx, a.db)
	if err != nil {
		return models.Storage{}, err
	}

	envelopeCount, err := repo.EnvelopeCount(ctx, a.db)
	if err != nil {
		return models.Storage{}, err
	}

	attachmentSize, err := a.fileStore.Size(ctx)
	if err != nil {
		return models.Storage{}, err
	}

	databaseSize, err := repo.Size(ctx, a.db)
	if err != nil {
		return models.Storage{}, err
	}

	return models.Storage{
		AttachmentCount: attachmentCount,
		EnvelopeCount:   envelopeCount,
		AttachmentSize:  attachmentSize,
		DatabaseSize:    databaseSize,
	}, nil
}

// - Tracer

func (a App) Tracer(source string) trace.Tracer {
	return trace.NewTracer(repo.NewTraceStore(a.db), source)
}

func (a App) TraceList(ctx context.Context, page pagination.Page, req models.DTOTraceListRequest) (models.DTOTraceListResult, error) {
	return repo.TraceList(ctx, a.db, page, req)
}

func (a App) TraceDrop(ctx context.Context) error {
	_, err := repo.TraceDrop(ctx, a.db)
	return err
}

// - RetentionPolicy

func (a App) RetentionPolicyGet(ctx context.Context) models.ConfigRetentionPolicy {
	return a.config.RetentionPolicy
}

func (a App) RetentionPolicyRun(ctx context.Context, tracer trace.Tracer) error {
	count1, err := retention.DeleteEnvelopeByAge(ctx, tracer, a.db, a.config.RetentionPolicy)
	if err != nil {
		return err
	}

	count2, err := retention.DeleteEnvelopeByCount(ctx, tracer, a.db, a.config.RetentionPolicy)
	if err != nil {
		return err
	}

	if count := count1 + count2; count > 0 {
		a.bus.EnvelopeDeleted(ctx)
	}

	_, err = retention.DeleteAttachmentBySize(ctx, tracer, a.fileStore, a.config.RetentionPolicy)
	if err != nil {
		return err
	}

	_, err = retention.DeleteTraceByAge(ctx, tracer, a.db, a.config.RetentionPolicy)
	if err != nil {
		return err
	}

	return nil
}

// - Mailman

func (a App) MailmanEnqueue(ctx context.Context, envelopeID int64) error {
	err := repo.MailmanEnqueue(ctx, a.db, envelopeID)
	if err != nil {
		return err
	}

	a.bus.MailmanEnqueued(ctx)

	return nil
}

func (a App) MailmanDequeue(ctx context.Context) (*models.Envelope, error) {
	envelopeID, err := repo.MailmanDequeue(ctx, a.db)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	env, err := repo.EnvelopeGet(ctx, a.db, envelopeID)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &env, nil
}

var _ core.App = App{}
