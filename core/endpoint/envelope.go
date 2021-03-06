package endpoint

import (
	"bytes"
	"context"
	"errors"
	"log"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

func (e Endpoint) Send(ctx context.Context, env *envelope.Envelope, atts []Attachment) error {
	// Text
	var text string
	if !e.TextDisable {
		var err error
		text, err = e.Text(env)
		if err != nil {
			return err
		}
	}

	// Attachments
	if e.AttachmentsDisable {
		atts = []Attachment{}
	}

	return e.SendRaw(ctx, text, atts)
}

func (e Endpoint) Text(env *envelope.Envelope) (string, error) {
	var buffer bytes.Buffer
	if err := e.textTemplate.Execute(&buffer, env); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func NewAttachment(att *envelope.Attachment, data []byte) Attachment {
	return Attachment{
		Name:    att.Name,
		Data:    data,
		IsImage: att.IsImage(),
	}
}

func convertAttachment(ctx context.Context, es envelope.Service, att envelope.Attachment) (Attachment, error) {
	data, err := es.GetData(ctx, &att)
	if err != nil {
		return Attachment{}, err
	}

	return NewAttachment(&att, data), nil
}

func ConvertAttachments(ctx context.Context, es envelope.Service, env *envelope.Envelope) ([]Attachment, error) {
	atts := []Attachment{}
	for _, att := range env.Attachments {
		newAtt, err := convertAttachment(ctx, es, att)
		if err != nil {
			if errors.Is(err, core.ErrDataNotFound) {
				log.Printf("endpoint.ConvertAttachments: name '%s': fileName '%s': %s", att.Name, att.FileName(), err)
				continue
			}

			return nil, err
		}

		atts = append(atts, newAtt)
	}

	return atts, nil
}
