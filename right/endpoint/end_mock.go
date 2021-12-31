package endpoint

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Mock struct{}

func NewMock() (*Mock, error) {
	return &Mock{}, nil
}

func (m *Mock) Send(message *domain.EndpointMessage) error {
	fmt.Printf("endpoint.Mock.Send: %+v\n", message)
	return nil
}