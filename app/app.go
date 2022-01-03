package app

import "github.com/ItsNotGoodName/smtpbridge/core"

type App struct {
	attachmentREPO core.AttachmentRepositoryPort
	messageREPO    core.MessageRepositoryPort
	endpointREPO   core.EndpointRepositoryPort
	authSVC        core.AuthServicePort
	bridgeSVC      core.BridgeServicePort
	endpointSVC    core.EndpointServicePort
	messageSVC     core.MessageServicePort
}

func New(
	attachmentREPO core.AttachmentRepositoryPort,
	messageREPO core.MessageRepositoryPort,
	endpointREPO core.EndpointRepositoryPort,
	authSVC core.AuthServicePort,
	bridgeSVC core.BridgeServicePort,
	endpointSVC core.EndpointServicePort,
	messageSVC core.MessageServicePort,
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
