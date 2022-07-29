package boltdb

import (
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type messageModel struct {
	ID        int64               `json:"id" storm:"id,increment"`
	CreatedAt time.Time           `json:"created_at"`
	Date      time.Time           `json:"date"`
	Subject   string              `json:"subject"`
	From      string              `json:"from"`
	To        map[string]struct{} `json:"to"`
	Text      string              `json:"text"`
	HTML      string              `json:"html"`
}

func messageMC(msg *messageModel) *envelope.Message {
	return &envelope.Message{
		ID:        msg.ID,
		CreatedAt: msg.CreatedAt,
		Date:      msg.Date,
		Subject:   msg.Subject,
		From:      msg.From,
		To:        msg.To,
		Text:      msg.Text,
		HTML:      msg.HTML,
	}
}

func messageCM(msg *envelope.Message) *messageModel {
	return &messageModel{
		ID:        msg.ID,
		CreatedAt: msg.CreatedAt,
		Date:      msg.Date,
		Subject:   msg.Subject,
		From:      msg.From,
		To:        msg.To,
		Text:      msg.Text,
		HTML:      msg.HTML,
	}
}

type attachmentModel struct {
	ID        int64  `json:"id" storm:"id,increment"`
	MessageID int64  `json:"message_id" storm:"index"`
	Name      string `json:"name"`
	Mime      string `json:"mime"`
	Extension string `json:"extension"`
}

func attachmentMC(attM *attachmentModel) *envelope.Attachment {
	return &envelope.Attachment{
		ID:        attM.ID,
		MessageID: attM.MessageID,
		Name:      attM.Name,
		Mime:      attM.Mime,
		Extension: attM.Extension,
	}
}

func attachmentCM(att *envelope.Attachment) *attachmentModel {
	return &attachmentModel{
		ID:        att.ID,
		MessageID: att.MessageID,
		Name:      att.Name,
		Mime:      att.Mime,
		Extension: att.Extension,
	}
}
