package app

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Info struct {
	AttachmentSize    int64
	AttachmentSizeMiB string
	MessageCount      int
	AttachmentCount   int
}

func (a *App) Info() (*Info, error) {
	msgCount, err := a.messageREPO.Count(&core.MessageParam{})
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

	return &Info{
		AttachmentSize:    attSize,
		MessageCount:      msgCount,
		AttachmentCount:   attCount,
		AttachmentSizeMiB: attSizeMiB,
	}, nil
}
