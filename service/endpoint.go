package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Endpoint struct {
	endpointREPO core.EndpointRepositoryPort
}

func NewEndpoint(endpointREPO core.EndpointRepositoryPort) *Endpoint {
	return &Endpoint{endpointREPO: endpointREPO}
}

func (e *Endpoint) SendBridges(msg *core.Message, bridges []*core.Bridge) (core.Status, error) {
	// TODO: refactor entire method
	if len(bridges) == 0 {
		log.Println("app.messageSend: no valid bridges: skipped message", msg.UUID)
		return core.StatusSkipped, nil
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
			var endpoint core.EndpointPort
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
		return core.StatusSent, errGet
	}

	if skipped > 0 {
		log.Println("app.messageSend: only_* produced empty message: skipped message", msg.UUID)
		return core.StatusSkipped, errGet
	}

	return core.StatusFailed, errGet
}
