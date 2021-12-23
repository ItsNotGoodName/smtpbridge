package endpoint

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Mock struct{}

func NewMock() (*Mock, error) {
	return &Mock{}, nil
}

func (m *Mock) Send(message *app.Message) error {
	fmt.Println("endpoint.Mock.Send:", message)
	return nil
}
