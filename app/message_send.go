package app

type MessageSendRequest struct {
	UUID string
}

func (a *App) MessageSend(req *MessageSendRequest) error {
	msg, err := a.messageREPO.Get(req.UUID)
	if err != nil {
		return err
	}

	atts, err := a.attachmentREPO.ListByMessage(msg)
	if err != nil {
		return err
	}

	for i := range atts {
		if err := a.attachmentREPO.LoadData(&atts[i]); err != nil {
			return err
		}
	}

	return a.endpointSVC.Process(msg, atts, a.bridgeSVC.ListByMessage(msg))
}
