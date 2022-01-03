package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Endpoint struct {
	endpointREPO core.EndpointRepositoryPort
}

func NewEndpoint(endpointREPO core.EndpointRepositoryPort) *Endpoint {
	e := Endpoint{
		endpointREPO: endpointREPO,
	}
	return &e
}

func (e *Endpoint) SendByEndpointNames(emsg *core.EndpointMessage, endpointNames []string) error {
	endpoints := make([]core.EndpointPort, len(endpointNames))
	for i, endpointName := range endpointNames {
		endpoint, err := e.endpointREPO.Get(endpointName)
		if err != nil {
			return err
		}

		endpoints[i] = endpoint
	}

	sent := false
	for _, endpoint := range endpoints {
		err := endpoint.Send(emsg)
		if err != nil {
			log.Println("service.Endpoint.SendByBridge:", err)
		} else {
			sent = true
		}
	}

	if !sent {
		return core.ErrEndpointSendFailed
	}

	return nil
}
