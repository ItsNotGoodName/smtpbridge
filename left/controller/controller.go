package controller

import (
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type Controller struct {
	envelopeService envelope.Service
}

func New(envelopeService envelope.Service) *Controller {
	return &Controller{
		envelopeService: envelopeService,
	}
}
