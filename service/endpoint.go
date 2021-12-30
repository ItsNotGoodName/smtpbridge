package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Endpoint struct {
	dao app.DAO
}

func NewEndpoint(dao app.DAO) *Endpoint {
	return &Endpoint{dao}
}

func (e *Endpoint) SendBridges(msg *app.Message, bridges []app.Bridge) error {
	if len(bridges) == 0 {
		return app.ErrBridgesNotFound
	}

	sent := 0
	for _, bridge := range bridges {
		emsg := bridge.EndpointMessage(msg)
		if emsg.IsEmpty() {
			continue
		}

		for _, name := range bridge.Endpoints {
			endpoint, err := e.dao.Endpoint.Get(name)
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
		return app.ErrEndpointSendFailed
	}

	return nil
}
