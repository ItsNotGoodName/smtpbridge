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
	// TODO: refactor entire method
	if len(bridges) == 0 {
		log.Println("app.messageSend: no valid bridges: skipped message", msg.UUID)
		return domain.StatusSkipped, nil
	}

	var errGet error
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
			endpoint, errGet = e.endpointREPO.Get(name)
			if errGet != nil {
				break
			}

			// TODO: worker pool
			if errEnd := endpoint.Send(emsg); errEnd != nil {
				log.Println("service.Endpoint.SendBridges:", errEnd)
			} else {
				sent++
			}
		}
	}

	if sent > 0 {
		log.Println("app.messageSend: sent message", msg.UUID)
		return domain.StatusSent, errGet
	}

	if skipped > 0 {
		log.Println("app.messageSend: only_* produced empty message: skipped message", msg.UUID)
		return domain.StatusSkipped, errGet
	}

	return domain.StatusFailed, errGet
}
