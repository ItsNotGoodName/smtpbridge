package janitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/cursor"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type JanitorService struct {
	attachmentRepository attachment.Repository
	dataRepository       attachment.DataRepository
	messageRepository    message.Repository
	maxSize              int64
}

func NewJanitorService(attachmentRepository attachment.Repository, messageRepository message.Repository, dataRepository attachment.DataRepository, maxSize int64) *JanitorService {
	return &JanitorService{
		attachmentRepository: attachmentRepository,
		dataRepository:       dataRepository,
		messageRepository:    messageRepository,
		maxSize:              maxSize,
	}
}

func (js *JanitorService) Clean(ctx context.Context) error {
	for {
		size, err := js.dataRepository.Size(ctx)
		if err != nil {
			return err
		}
		if size < js.maxSize {
			return nil
		}

		listParam := message.ListParam{
			Cursor: cursor.NewOldest(5),
		}
		if err = js.messageRepository.List(ctx, &listParam); err != nil {
			return err
		}
		msgs := listParam.Messages

		if len(msgs) == 0 {
			return fmt.Errorf("attachments are over capacity (%d bytes > %d bytes), but no messages available to delete", size, js.maxSize)
		}

		for i := range msgs {
			err := js.messageRepository.Delete(ctx, &msgs[i])
			if err != nil {
				return err
			}
			attsCount, err := js.attachmentRepository.CountByMessage(ctx, &msgs[i])
			if err != nil {
				log.Println("janitor.JanitorService.Clean: could not count attachments:", err)
			}
			log.Printf("janitor.JanitorService.Clean: deleted message '%d' with %d attachments", msgs[i].ID, attsCount)
		}
	}
}

func (j *JanitorService) Run(ctx context.Context, done chan<- struct{}) {
	log.Println("janitor.JanitorService.Run: started")

	if j.maxSize <= 0 {
		log.Printf("janitor.JanitorService.Run: maxSize is %d, janitor is disabled", j.maxSize)
		done <- struct{}{}
		return
	}

	clean := func() {
		if err := j.Clean(ctx); err != nil {
			log.Printf("janitor.JanitorService.Run: %s", err)
		}
	}
	clean()

	t := time.NewTicker(time.Minute * 10)

	for {
		select {
		case <-t.C:
			clean()
		case <-ctx.Done():
			log.Println("janitor.JanitorService.Run: stopped")
			done <- struct{}{}
			return
		}
	}
}
