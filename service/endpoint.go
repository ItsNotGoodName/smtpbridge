package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Endpoint struct {
	endpointREPO domain.EndpointRepositoryPort
}

func NewEndpoint(endpointREPO domain.EndpointRepositoryPort) *Endpoint {
	return &Endpoint{endpointREPO: endpointREPO}
}

func (e *Endpoint) SendBridges(msg *domain.Message, bridges []domain.Bridge) error {
	if len(bridges) == 0 {
		return domain.ErrBridgesNotFound
	}

	sent := 0
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

			// TODO: worker pool
			if err = endpoint.Send(emsg); err != nil {
				log.Println("service.Endpoint.SendBridges:", err)
			} else {
				sent++
			}
		}
	}

	if sent == 0 {
		return domain.ErrEndpointSendFailed
	}

	return nil
}
