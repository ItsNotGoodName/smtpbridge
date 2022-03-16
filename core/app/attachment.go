package app

import (
	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/dto"
)

func newAttachment(att *attachment.Attachment) dto.Attachment {
	return dto.Attachment{
		ID:   att.ID,
		Name: att.Name,
		File: att.File(),
		Type: string(att.Type),
	}
}
