package event

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

type EnvelopeService struct {
	envelopeService envelope.Service
	pub             *Pub
}

func NewEnvelopeService(envelopeService envelope.Service, pub *Pub) *EnvelopeService {
	return &EnvelopeService{
		envelopeService: envelopeService,
		pub:             pub,
	}
}

func (e *EnvelopeService) ListAttachment(ctx context.Context, page *paginate.Page) ([]envelope.Attachment, error) {
	return e.envelopeService.ListAttachment(ctx, page)
}

func (e *EnvelopeService) ListEnvelope(ctx context.Context, page *paginate.Page) ([]envelope.Envelope, error) {
	return e.envelopeService.ListEnvelope(ctx, page)
}

func (e *EnvelopeService) GetEnvelope(ctx context.Context, msgID int64) (*envelope.Envelope, error) {
	return e.envelopeService.GetEnvelope(ctx, msgID)
}

func (e *EnvelopeService) CreateEnvelope(ctx context.Context, req *envelope.CreateEnvelopeRequest) (int64, error) {
	id, err := e.envelopeService.CreateEnvelope(ctx, req)
	if err != nil {
		return 0, err
	}

	env, err := e.GetEnvelope(ctx, id)
	if err != nil {
		return id, err
	}

	e.pub.Publish(Event{TopicEnvelopeCreated, env})

	return id, nil
}

func (e *EnvelopeService) DeleteEnvelope(ctx context.Context, msgID int64) error {
	return e.envelopeService.DeleteEnvelope(ctx, msgID)
}

func (e *EnvelopeService) GetData(ctx context.Context, att *envelope.Attachment) ([]byte, error) {
	return e.envelopeService.GetData(ctx, att)
}
