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

func (e *Endpoint) SendBridges(msg *domain.Message, bridges []*domain.Bridge) (domain.Status, error) {
	if len(bridges) == 0 {
		return domain.StatusSkipped, nil
	}

	var err error
	sent := 0
	skipped := 0
	for _, bridge := range bridges {
		emsg := bridge.EndpointMessage(msg)
		if emsg.IsEmpty() {
			skipped++
			continue
		}

		for _, name := range bridge.Endpoints {
			var endpoint domain.EndpointPort
			endpoint, err = e.endpointREPO.Get(name)
			if err != nil {
				break
			}

			// TODO: worker pool
			if err = endpoint.Send(emsg); err != nil {
				log.Println("service.Endpoint.SendBridges:", err)
			} else {
				sent++
			}
		}
	}

	if sent > 0 {
		return domain.StatusSent, err
	}

	if skipped > 0 {
		return domain.StatusSkipped, err
	}

	return domain.StatusFailed, err
}
