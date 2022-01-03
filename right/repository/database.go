package repository

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/asdine/storm"
)

type Database struct {
	db         *storm.DB
	attachment *Attachment
	message    *Message
}

func NewDatabase(cfg *config.Config) Database {
	db, err := storm.Open(cfg.DB.DB)
	if err != nil {
		log.Fatalln("repository.NewStorm: could not open database:", err)
	}

	attachment := NewAttachment(cfg, db)
	message := NewMessage(db, attachment)

	return Database{
		db:         db,
		attachment: attachment,
		message:    message,
	}
}

func (d Database) Close() error {
	return d.db.Close()
}

func (d Database) AttachmentRepository() core.AttachmentRepositoryPort {
	return d.attachment
}

func (d Database) MessageRepository() core.MessageRepositoryPort {
	return d.message
}
