package procs

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/db"
	"github.com/ItsNotGoodName/smtpbridge/internal/files"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

func StorageGet(cc core.Context) (models.Storage, error) {
	attachmentCount, err := db.EnvelopeAttachmentCount(cc)
	if err != nil {
		return models.Storage{}, err
	}
	envelopeCount, err := db.EnvelopeCount(cc)
	if err != nil {
		return models.Storage{}, err
	}
	attachmentSize, err := files.Size(cc)
	if err != nil {
		return models.Storage{}, err
	}
	databaseSize, err := db.Size(cc)
	if err != nil {
		return models.Storage{}, err
	}

	storage := models.Storage{
		AttachmentCount: attachmentCount,
		EnvelopeCount:   envelopeCount,
		AttachmentSize:  attachmentSize,
		DatabaseSize:    databaseSize,
	}

	return storage, nil
}
