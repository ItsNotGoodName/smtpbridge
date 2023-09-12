package senders

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

type Console struct{}

func NewConsole() Console {
	return Console{}
}

func (c Console) Send(ctx context.Context, env models.Envelope, tr Transformer) error {
	title, err := tr.Title(ctx, env)
	if err != nil {
		return err
	}

	body, err := tr.Body(ctx, env)
	if err != nil {
		return err
	}

	if title == "" && body == "" {
		return nil
	}

	fmt.Printf("Title: %s\nBody: %s\n", title, body)

	return nil
}
