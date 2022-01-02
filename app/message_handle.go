package app

import "log"

type MessageHandleRequest struct {
	Subject               string
	From                  string
	To                    map[string]struct{}
	Text                  string
	IgnoreAttachmentError bool
	attachments           []attachmentHandleRequest
}

type attachmentHandleRequest struct {
	name string
	data []byte
}

func (c *MessageHandleRequest) AddAttachment(name string, data []byte) {
	c.attachments = append(c.attachments, attachmentHandleRequest{name, data})
}

func (a *App) MessageHandle(req *MessageHandleRequest) error {
	msg, err := a.messageSVC.Create(req.Subject, req.From, req.To, req.Text)
	if err != nil {
		return err
	}

	for _, attachment := range req.attachments {
		_, err = a.messageSVC.CreateAttachment(msg, attachment.name, attachment.data)
		if err != nil {
			if !req.IgnoreAttachmentError {
				return err
			}
			log.Println("app.App.MessageHandle: could not create attachment:", err)
		}
	}

	// TODO: send to service
	return a.messageSend(&messageSendRequest{
		Message: msg,
	})
}
