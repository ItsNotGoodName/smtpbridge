package app

type MessageGetRequest struct {
	UUID string
}

type MessageGetResponse struct {
	Message Message
}

func (a *App) MessageGet(req *MessageGetRequest) (*MessageGetResponse, error) {
	msg, err := a.messageREPO.Get(req.UUID)
	if err != nil {
		return nil, err
	}

	atts, err := a.attachmentREPO.ListByMessage(msg)
	if err != nil {
		return nil, err
	}

	return &MessageGetResponse{Message: NewMessage(msg, atts)}, nil
}
