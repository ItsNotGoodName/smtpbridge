package endpoints

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
)

type Mock struct{}

func NewMock() (*Mock, error) {
	return &Mock{}, nil
}

type mockJSON struct {
	Title             string `json:"title"`
	Body              string `json:"body"`
	AttachmentsLength int    `json:"attachments_length"`
}

func (m *Mock) Send(ctx context.Context, env endpoint.Envelope) error {
	mj := mockJSON{
		Title:             env.Message.Title,
		Body:              env.Message.Body,
		AttachmentsLength: len(env.Attachments),
	}

	mb, err := json.Marshal(mj)
	if err != nil {
		return err
	}

	fmt.Println(string(mb))

	return nil
}
