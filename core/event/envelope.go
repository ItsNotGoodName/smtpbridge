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

func (es *EnvelopeService) ListAttachment(ctx context.Context, page *paginate.Page) ([]envelope.Attachment, error) {
	return es.envelopeService.ListAttachment(ctx, page)
}

func (es *EnvelopeService) ListEnvelope(ctx context.Context, page *paginate.Page) ([]envelope.Envelope, error) {
	return es.envelopeService.ListEnvelope(ctx, page)
}

func (es *EnvelopeService) GetEnvelope(ctx context.Context, msgID int64) (*envelope.Envelope, error) {
	return es.envelopeService.GetEnvelope(ctx, msgID)
}

func (es *EnvelopeService) CreateEnvelope(ctx context.Context, req *envelope.CreateEnvelopeRequest) (int64, error) {
	id, err := es.envelopeService.CreateEnvelope(ctx, req)
	if err != nil {
		return 0, err
	}

	env, err := es.GetEnvelope(ctx, id)
	if err != nil {
		return id, err
	}

	es.pub.Publish(Event{TopicEnvelopeCreated, env})

	return id, nil
}

func (es *EnvelopeService) DeleteEnvelope(ctx context.Context, msgID int64) error {
	return es.envelopeService.DeleteEnvelope(ctx, msgID)
}

func (es *EnvelopeService) GetData(ctx context.Context, att *envelope.Attachment) ([]byte, error) {
	return es.envelopeService.GetData(ctx, att)
}

func (es *EnvelopeService) GetDataSize(ctx context.Context) (int64, error) {
	return es.envelopeService.GetDataSize(ctx)
}
