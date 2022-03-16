package mockdb

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type Message struct{}

func NewMessage() *Message {
	return &Message{}
}

func (Message) Create(ctx context.Context, msg *message.Message) error {
	return nil
}

func (Message) Count(ctx context.Context) (int, error) {
	return 0, nil
}

func (Message) Get(ctx context.Context, id int64) (*message.Message, error) {
	return nil, message.ErrNotFound
}

func (Message) List(ctx context.Context, param *message.ListParam) error {
	return nil
}

func (Message) Update(ctx context.Context, msg *message.Message, updateFN func(msg *message.Message) (*message.Message, error)) error {
	return nil
}

func (Message) Delete(ctx context.Context, msg *message.Message) error {
	return nil
}
