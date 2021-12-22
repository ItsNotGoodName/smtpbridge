package endpoint

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Mock struct{}

func NewMock() (*Mock, error) {
	return &Mock{}, nil
}

func (m *Mock) Send(message *app.Message) error {
	log.Println("Mock:", message)
	return nil
}
