package app

import "github.com/ItsNotGoodName/smtpbridge/domain"

type App struct {
	dao         domain.DAO
	authSVC     domain.AuthServicePort
	bridgeSVC   domain.BridgeServicePort
	endpointSVC domain.EndpointServicePort
	messageSVC  domain.MessageServicePort
}

func New(
	dao domain.DAO,
	authSVC domain.AuthServicePort,
	bridgeSVC domain.BridgeServicePort,
	endpointSVC domain.EndpointServicePort,
	messageSVC domain.MessageServicePort,
) *App {
	return &App{
		dao,
		authSVC,
		bridgeSVC,
		endpointSVC,
		messageSVC,
	}
}
