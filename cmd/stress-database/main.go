package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/internal/repo"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, err := database.New("./smtpbridge_data/smtpbridge.db", false)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		msg := models.Message{}
		for {
			id, err := repo.EnvelopeCreate(ctx, db, msg, []models.Attachment{})
			if err != nil {
				log.Fatalln(err)
			}

			err = repo.MailmanEnqueue(ctx, db, id)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}()

	for {
		id, err := repo.MailmanDequeue(ctx, db)
		if err != nil {
			if !errors.Is(err, repo.ErrNoRows) {
				log.Fatalln(err)
			}

			time.Sleep(1 * time.Second)
		}

		fmt.Println(id)
	}
}
