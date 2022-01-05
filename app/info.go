package app

import "fmt"

type InfoResponse struct {
	AttachmentSize    int64
	AttachmentSizeMiB string
	MessageCount      int
	AttachmentCount   int
}

func (a *App) Info() (*InfoResponse, error) {
	msgCount, err := a.messageREPO.Count()
	if err != nil {
		return nil, err
	}

	attCount, err := a.attachmentREPO.Count()
	if err != nil {
		return nil, err
	}

	attSize, err := a.attachmentREPO.GetSizeAll()
	if err != nil {
		return nil, err
	}

	attSizeMiB := fmt.Sprintf("%.2f", float64(attSize)/1024/1024)

	return &InfoResponse{
		AttachmentSize:    attSize,
		MessageCount:      msgCount,
		AttachmentCount:   attCount,
		AttachmentSizeMiB: attSizeMiB,
	}, nil
}
