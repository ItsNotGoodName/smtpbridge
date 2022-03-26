package app

import (
	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/auth"
	"github.com/ItsNotGoodName/smtpbridge/core/bridge"
	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type App struct {
	attachmentRepository     attachment.Repository
	attachmentService        attachment.Service
	bridgeService            bridge.Service
	attachmentDataRepository attachment.DataRepository
	endpointService          endpoint.Service
	eventRepository          event.Repository
	messageRepository        message.Repository
	messageService           message.Service
	smtpAuthService          auth.Service
}

func New(
	attachmentRepository attachment.Repository,
	attachmentService attachment.Service,
	bridgeService bridge.Service,
	attachmentDataRepository attachment.DataRepository,
	endpointService endpoint.Service,
	eventRepository event.Repository,
	messageRepository message.Repository,
	messageService message.Service,
	smtpAuthService auth.Service,
) *App {
	return &App{
		attachmentRepository,
		attachmentService,
		bridgeService,
		attachmentDataRepository,
		endpointService,
		eventRepository,
		messageRepository,
		messageService,
		smtpAuthService,
	}
}

func newMessage(msg *message.Message, atts []attachment.Attachment) dto.Message {
	var attachments []dto.Attachment = make([]dto.Attachment, 0, len(atts))
	for _, att := range atts {
		attachments = append(attachments, newAttachment(&att))
	}

	var to []string
	for toAddr := range msg.To {
		to = append(to, toAddr)
	}

	return dto.Message{
		ID:          msg.ID,
		CreatedAt:   msg.CreatedAt,
		From:        msg.From,
		To:          to,
		Subject:     msg.Subject,
		Text:        msg.Text,
		Attachments: attachments,
	}
}

func newAttachment(att *attachment.Attachment) dto.Attachment {
	return dto.Attachment{
		ID:   att.ID,
		Name: att.Name,
		File: att.File(),
		Type: string(att.Type),
	}
}

func newEvent(e *event.Event) *dto.Event {
	return &dto.Event{
		ID:          e.ID,
		Name:        string(e.Name),
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
	}
}
