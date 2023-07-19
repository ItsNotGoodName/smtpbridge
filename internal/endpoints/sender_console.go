package endpoints

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
)

type Console struct{}

func NewConsole() Console {
	return Console{}
}

func (c Console) Send(cc *core.Context, env envelope.Envelope, config Config) error {
	body, err := GetBody(env, config)
	if err != nil {
		return err
	}
	if body == "" {
		return nil
	}

	fmt.Println(body)

	return nil
}
