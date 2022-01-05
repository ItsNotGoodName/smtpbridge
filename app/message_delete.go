package app

type MessageDeleteRequest struct {
	UUID string
}

func (a *App) MessageDelete(req MessageDeleteRequest) error {
	msg, err := a.messageREPO.Get(req.UUID)
	if err != nil {
		return err
	}

	return a.messageREPO.Delete(msg)
}
