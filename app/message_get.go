package app

type MessageGetRequest struct {
	UUID string
}

type MessageGetResponse struct {
	Message Message
}

func (a *App) MessageGet(req *MessageGetRequest) (*MessageGetResponse, error) {
	msg, err := a.messageSVC.Get(req.UUID)
	if err != nil {
		return nil, err
	}

	return &MessageGetResponse{Message: NewMessage(msg)}, nil
}
