package app

type MessageSendRequest struct {
	UUID string
}

func (a *App) MessageSend(req *MessageSendRequest) error {
	msg, err := a.messageSVC.Get(req.UUID)
	if err != nil {
		return err
	}

	if err := a.messageSVC.LoadData(msg); err != nil {
		return err
	}

	return a.endpointSVC.Process(msg, a.bridgeSVC.ListByMessage(msg))
}
