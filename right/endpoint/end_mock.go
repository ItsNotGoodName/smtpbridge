package endpoint

import (
	"encoding/json"
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Mock struct{}

func NewMock() (*Mock, error) {
	return &Mock{}, nil
}

type MockJSON struct {
	Text              string `json:"text"`
	AttachmentsLength int    `json:"attachments_length"`
}

func (m *Mock) Send(message *core.EndpointMessage) error {
	mj := MockJSON{
		Text:              message.Text,
		AttachmentsLength: len(message.Attachments),
	}

	mb, err := json.Marshal(mj)
	if err != nil {
		return err
	}

	fmt.Println(string(mb))
	return nil
}
