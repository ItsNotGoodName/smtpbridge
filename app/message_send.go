package app

type MessageSendRequest struct {
	UUID string
}

func (a *App) MessageSend(req *MessageSendRequest) error {
	msg, err := a.messageSVC.Get(req.UUID)
	if err != nil {
		return err
	}

	for i := range msg.Attachments {
		var err error
		msg.Attachments[i].Data, err = a.attachmentREPO.GetData(&msg.Attachments[i])
		if err != nil {
			return err
		}
	}

	return a.messageSVC.Process(msg, a.bridgeSVC.ListByMessage(msg))
}
