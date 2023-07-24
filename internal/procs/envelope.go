package procs

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
	"github.com/ItsNotGoodName/smtpbridge/internal/events"
	"github.com/ItsNotGoodName/smtpbridge/internal/files"
	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
)

func EnvelopeMessageList(cc core.Context, page pagination.Page, filter envelope.MessageFilter) (envelope.MessageListResult, error) {
	return db.EnvelopeMessageList(cc, page, filter)
}

func EnvelopeCreate(cc core.Context, msg *envelope.Message, datts []envelope.DataAttachment) (int64, error) {
	// Extract attachments
	atts := make([]*envelope.Attachment, len(datts))
	for i := range atts {
		atts[i] = datts[i].Attachment
	}

	// Save message and attachments
	msgID, _, err := db.EnvelopeCreate(cc, msg, atts)
	if err != nil {
		return 0, err
	}

	// Save attachment's data
	for _, datt := range datts {
		if err := files.CreateFile(cc, datt); err != nil {
			return 0, err
		}
	}

	// Publish event
	events.PublishEnvelopeCreated(cc, msgID)

	return msgID, nil
}

func EnvelopeGet(cc core.Context, id int64) (envelope.Envelope, error) {
	return db.EnvelopeGet(cc, id)
}

func EnvelopeMessageHTMLGet(cc core.Context, id int64) (string, error) {
	return db.EnvelopeMessageHTMLGet(cc, id)
}

func EnvelopeDelete(cc core.Context, id int64) error {
	err := db.EnvelopeDelete(cc, id)
	if err != nil {
		return err
	}

	events.PublishEnvelopeDeleted(cc, id)

	return nil
}

func EnvelopeDeleteAll(cc core.Context) error {
	err := db.EnvelopeDeleteAll(cc)
	if err != nil {
		return err
	}

	events.PublishEnvelopeDeleted(cc)

	return nil
}

func EnvelopeAttachmentList(cc core.Context, page pagination.Page, filter envelope.AttachmentFilter) (envelope.AttachmentListResult, error) {
	return db.EnvelopeAttachmentList(cc, page, filter)
}
