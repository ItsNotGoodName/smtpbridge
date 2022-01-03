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

	errC := make(chan error, len(endpoints))
	for _, end := range endpoints {
		go func(emessage *core.EndpointMessage, endpoint core.EndpointPort) {
			errC <- endpoint.Send(emessage)
		}(emsg, end)
	}

	sent := false
	for i := 0; i < len(endpoints); i++ {
		err := <-errC
		if err != nil {
			log.Println("service.Endpoint.SendByEndpointNames:", err)
		} else {
			sent = true
		}
	}

	if !sent {
		return core.ErrEndpointSendFailed
	}

	return nil
}
