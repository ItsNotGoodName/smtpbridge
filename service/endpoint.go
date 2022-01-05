package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Endpoint struct {
	endpointREPO core.EndpointRepositoryPort
	messageSVC   core.MessageServicePort
}

func NewEndpoint(endpointREPO core.EndpointRepositoryPort, messageSVC core.MessageServicePort) *Endpoint {
	e := Endpoint{
		endpointREPO: endpointREPO,
		messageSVC:   messageSVC,
	}
	return &e
}

func (e *Endpoint) sendByEndpointNames(emsg *core.EndpointMessage, endpointNames []string) error {
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

func (e *Endpoint) Process(msg *core.Message, atts []core.Attachment, bridges []*core.Bridge) error {
	skipped := 0
	failed := 0
	for _, bridge := range bridges {
		emsg := bridge.EndpointMessage(msg, atts)
		if emsg.IsEmpty() {
			skipped++
			continue
		}

		if err := e.sendByEndpointNames(emsg, bridge.Endpoints); err != nil {
			failed++
			log.Println("service.Endpoint.Process:", err)
		}
	}

	length := len(bridges)
	status := core.StatusSent
	if skipped == length {
		status = core.StatusSkipped
	} else if failed+skipped == len(bridges) {
		status = core.StatusFailed
	}

	return e.messageSVC.UpdateStatus(msg, status)
}
