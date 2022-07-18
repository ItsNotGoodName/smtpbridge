package view

import (
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

const (
	ErrorPage string = "error.html"
	IndexPage string = "index.html"
)

type IndexData struct {
	Envelopes []envelope.Envelope
	Page      paginate.Page
}
