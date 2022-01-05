package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core"
)

type Janitor struct {
	messageREPO    core.MessageRepositoryPort
	attachmentREPO core.AttachmentRepositoryPort
	size           int64
	disabled       bool
}

func NewJanitor(cfg *config.Config, attachmentREPO core.AttachmentRepositoryPort, messageREPO core.MessageRepositoryPort) *Janitor {
	return &Janitor{
		attachmentREPO: attachmentREPO,
		messageREPO:    messageREPO,
		size:           cfg.DB.Size,
		disabled:       cfg.DB.Size == 0,
	}
}

func (j *Janitor) clean() error {
	for {
		attSize, err := j.attachmentREPO.GetSizeAll()
		if err != nil {
			return err
		}
		msgSize, err := j.messageREPO.GetSizeAll()
		if err != nil {
			return err
		}
		size := attSize + msgSize

		if size < j.size {
			return nil
		}

		msgs, err := j.messageREPO.ListOldest(5)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			return fmt.Errorf("database is over capacity (%d bytes > %d bytes), but no messages to delete", size, j.size)
		}

		for i := range msgs {
			err := j.messageREPO.Delete(&msgs[i])
			if err != nil {
				return err
			}
			attsCount, err := j.attachmentREPO.CountByMessage(&msgs[i])
			if err != nil {
				log.Println("service.Janitor.clean: could not count attachments:", err)
			}
			log.Printf("service.Janitor.clean: deleted message '%s' with %d attachments", msgs[i].UUID, attsCount)
		}
	}
}

func (j *Janitor) Run(ctx context.Context, done chan<- struct{}) {
	if j.disabled {
		log.Printf("service.Janitor.Run: disabled, database max size is %d bytes", j.size)
		done <- struct{}{}
		return
	}

	log.Println("service.Janitor.Run: started")

	clean := func() {
		if err := j.clean(); err != nil {
			log.Printf("service.Janitor.Run: %s", err)
		}
	}
	clean()

	t := time.NewTicker(time.Minute * 10)

	for {
		select {
		case <-t.C:
			clean()
		case <-ctx.Done():
			log.Println("service.Janitor.Run: stopped")
			done <- struct{}{}
			return
		}
	}
}
