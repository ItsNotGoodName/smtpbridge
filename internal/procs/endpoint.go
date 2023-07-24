package procs

import (
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
	"github.com/ItsNotGoodName/smtpbridge/internal/endpoints"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
)

func EndpointSend(cc core.Context, envelope_id int64, endpoint_id int64) error {
	end, err := db.EndpointGet(cc, endpoint_id)
	if err != nil {
		return err
	}

	env, err := db.EnvelopeGet(cc, envelope_id)
	if err != nil {
		return err
	}

	parsedEnd, err := end.Parse()
	if err != nil {
		return err
	}

	return parsedEnd.Send(cc, env)
}

func EndpointTest(cc core.Context, id int64) error {
	end, err := db.EndpointGet(cc, id)
	if err != nil {
		return err
	}

	parsedEnd, err := end.Parse()
	if err != nil {
		return err
	}

	msg := envelope.NewMessage(envelope.CreateMessage{
		Subject: "Test Subject",
		Text:    "Test Text",
		Date:    time.Now(),
	})
	env := envelope.Envelope{Message: msg}

	return parsedEnd.Send(cc, env)
}

func EndpointList(cc core.Context) ([]endpoints.Endpoint, error) {
	return db.EndpointList(cc)
}
