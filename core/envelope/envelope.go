package envelope

import (
	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type Envelope struct {
	Message     *message.Message
	Attachments []attachment.Attachment
}

func New(msg *message.Message, atts []attachment.Attachment) Envelope {
	return Envelope{
		Message:     msg,
		Attachments: atts,
	}
}
