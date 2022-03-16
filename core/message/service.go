package message

import "context"

type MessageService struct {
	messageRepository Repository
}

func NewMessageService(messageRepository Repository) *MessageService {
	return &MessageService{
		messageRepository: messageRepository,
	}
}

func (ms *MessageService) Create(ctx context.Context, param *Param) (*Message, error) {
	msg := New(param)

	return msg, ms.messageRepository.Create(ctx, msg)
}

func (ms *MessageService) Processed(ctx context.Context, msg *Message) error {
	return ms.messageRepository.Update(ctx, msg, func(msg *Message) (*Message, error) {
		msg.SetProcessed()
		return msg, nil
	})
}
