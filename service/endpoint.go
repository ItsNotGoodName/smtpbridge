package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Endpoint struct {
	bridgeSVC    app.BridgeServicePort
	messageSVC   app.MessageServicePort
	endpointREPO app.EndpointRepositoryPort
}

func NewEndpoint(bridgeSVC app.BridgeServicePort, messageSVC app.MessageServicePort, endpointREPO app.EndpointRepositoryPort) *Endpoint {
	return &Endpoint{
		bridgeSVC:    bridgeSVC,
		messageSVC:   messageSVC,
		endpointREPO: endpointREPO,
	}
}

func (e *Endpoint) Send(msg *app.Message) error {
	return e.SendBridges(msg, e.bridgeSVC.GetBridges(msg))
}

func (e *Endpoint) SendBridges(msg *app.Message, bridges []app.Bridge) error {
	if len(bridges) == 0 {
		return app.ErrBridgesNotFound
	}

	sentCount := 0
	for _, bridge := range bridges {
		emsg := bridge.EndpointMessage(msg)
		if emsg.IsEmpty() {
			continue
		}

		for _, name := range bridge.Endpoints {
			endpoint, err := e.endpointREPO.Get(name)
			if err != nil {
				return err
			}

			err = endpoint.Send(emsg)
			if err != nil {
				log.Println("service.Endpoint.SendBridges:", err)
			} else {
				sentCount++
			}
		}
	}

	if sentCount == 0 {
		if err := e.messageSVC.UpdateStatus(msg, app.StatusFailed); err != nil {
			return err
		}

		return app.ErrEndpointSendFailed
	}

	return e.messageSVC.UpdateStatus(msg, app.StatusSent)
}
