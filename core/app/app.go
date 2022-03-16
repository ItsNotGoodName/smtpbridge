package app

import (
	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/auth"
	"github.com/ItsNotGoodName/smtpbridge/core/bridge"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type App struct {
	attachmentRepository attachment.Repository
	attachmentService    attachment.Service
	bridgeService        bridge.Service
	dataRepository       attachment.DataRepository
	endpointService      endpoint.Service
	eventRepository      event.Repository
	messageRepository    message.Repository
	messageService       message.Service
	smtpAuthService      auth.Service
}

func New(
	attachmentRepository attachment.Repository,
	attachmentService attachment.Service,
	bridgeService bridge.Service,
	dataRepository attachment.DataRepository,
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
		dataRepository,
		endpointService,
		eventRepository,
		messageRepository,
		messageService,
		smtpAuthService,
	}
}
