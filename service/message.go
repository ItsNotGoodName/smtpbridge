package service

import (
	"log"

	"github.com/ItsNotGoodName/go-smtpbridge/app"
)

type Message struct{}

func NewMessage() *Message {
	return &Message{}
}

func (m *Message) Handle(msg *app.Message) error {
	log.Println("service.Message.Handle:", msg)
	return nil
}
