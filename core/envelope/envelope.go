package envelope

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

type (
	Envelope struct {
		Message     Message
		Attachments []Attachment
	}

	CreateEnvelopeRequest struct {
		From       string
		To         []string
		Subject    string
		Text       string
		HTML       string
		Attachment []CreateAttachmentRequest
	}

	CreateAttachmentRequest struct {
		Name string
		Data []byte
	}

	Service interface {
		ListEnvelope(ctx context.Context, page *paginate.Page) ([]Envelope, error)
		CreateEnvelope(ctx context.Context, req *CreateEnvelopeRequest) (int64, error)
		DeleteEnvelope(ctx context.Context, msgID int64) error
	}

	Store interface {
		ListEnvelope(ctx context.Context, offset, limit int, ascending bool) ([]Envelope, int, error)
		CreateMessage(ctx context.Context, msg *Message) error
		DeleteMessage(ctx context.Context, msgID int64) error
		GetMessage(ctx context.Context, id int64) (*Message, error)
		CreateAttachment(ctx context.Context, att *Attachment) error
		GetAttachment(ctx context.Context, id int64) (*Attachment, error)
		GetAttachmentsByMessageID(ctx context.Context, msgID int64) ([]Attachment, error)
	}

	DataStore interface {
		CreateData(ctx context.Context, att *Attachment, data []byte) error
		GetData(ctx context.Context, att *Attachment) ([]byte, error)
		DeleteData(ctx context.Context, att *Attachment) error
	}
)

type EnvelopeService struct {
	store     Store
	dataStore DataStore
}

func NewEnvelopeService(store Store, dataStore DataStore) *EnvelopeService {
	return &EnvelopeService{
		store:     store,
		dataStore: dataStore,
	}
}

func (es *EnvelopeService) ListEnvelope(ctx context.Context, page *paginate.Page) ([]Envelope, error) {
	envs, count, err := es.store.ListEnvelope(ctx, page.Offset(), page.Limit, page.Ascending)
	if err != nil {
		return nil, err
	}

	page.SetMaxCount(count)

	return envs, nil
}

func (es *EnvelopeService) CreateEnvelope(ctx context.Context, req *CreateEnvelopeRequest) (int64, error) {
	msg := NewMessage(req.From, req.To, req.Subject, req.Text, req.HTML)
	if err := es.store.CreateMessage(ctx, msg); err != nil {
		return 0, err
	}

	for _, attReq := range req.Attachment {
		att := NewAttachment(msg.ID, attReq.Name, attReq.Data)
		if err := es.store.CreateAttachment(ctx, att); err != nil {
			return 0, err
		}

		if err := es.dataStore.CreateData(ctx, att, attReq.Data); err != nil {
			return 0, err
		}
	}

	return msg.ID, nil
}

func (es *EnvelopeService) DeleteEnvelope(ctx context.Context, msgID int64) error {
	atts, err := es.store.GetAttachmentsByMessageID(ctx, msgID)
	if err != nil {
		return err
	}

	if err := es.store.DeleteMessage(ctx, msgID); err != nil {
		return err
	}

	for _, att := range atts {
		if err := es.dataStore.DeleteData(ctx, &att); err != nil {
			return err
		}
	}

	return nil
}
