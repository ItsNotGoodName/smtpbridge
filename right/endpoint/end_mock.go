package endpoint

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Mock struct{}

func NewMock() (*Mock, error) {
	return &Mock{}, nil
}

func (m *Mock) Send(message *core.EndpointMessage) error {
	fmt.Printf("endpoint.Mock.Send: %+v\n", message)
	return nil
}
