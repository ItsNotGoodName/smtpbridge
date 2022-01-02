package app

import "github.com/ItsNotGoodName/smtpbridge/domain"

type App struct {
	attachmentREPO domain.AttachmentRepositoryPort
	messageREPO    domain.MessageRepositoryPort
	endpointREPO   domain.EndpointRepositoryPort
	authSVC        domain.AuthServicePort
	bridgeSVC      domain.BridgeServicePort
	endpointSVC    domain.EndpointServicePort
	messageSVC     domain.MessageServicePort
}

func New(
	attachmentREPO domain.AttachmentRepositoryPort,
	messageREPO domain.MessageRepositoryPort,
	endpointREPO domain.EndpointRepositoryPort,
	authSVC domain.AuthServicePort,
	bridgeSVC domain.BridgeServicePort,
	endpointSVC domain.EndpointServicePort,
	messageSVC domain.MessageServicePort,
) *App {
	return &App{
		attachmentREPO,
		messageREPO,
		endpointREPO,
		authSVC,
		bridgeSVC,
		endpointSVC,
		messageSVC,
	}
}
