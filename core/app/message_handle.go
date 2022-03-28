package app

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

func (a *App) MessageHandle(ctx context.Context, req *dto.MessageHandleRequest) error {
	msg, err := a.messageService.Create(ctx, &message.Param{
		Subject: req.Subject,
		From:    req.From,
		To:      req.To,
		Text:    req.Text,
	})
	if err != nil {
		return err
	}

	var atts []attachment.Attachment
	for _, att := range req.Attachments {
		att, err := a.attachmentService.Create(ctx, &attachment.Param{
			Name:    att.Name,
			Data:    att.Data,
			Message: msg,
		})
		if err != nil {
			log.Println("app.App.MessageHandle: could not create attachment:", err)
			continue
		}

		atts = append(atts, *att)
	}

	env := envelope.New(msg, atts)
	return a.bridgeService.HandleEnvelope(ctx, a.bridgeService.ListByEnvelope(env), env)
}
