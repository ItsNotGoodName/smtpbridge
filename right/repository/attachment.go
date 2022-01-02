package repository

import (
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/domain"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type Attachment struct {
	attDir string
	fs     fs.FS
	db     *storm.DB
}

func NewAttachment(cfg *config.Config, db *storm.DB) *Attachment {
	err := os.MkdirAll(cfg.DB.Attachments, 0755)
	if err != nil {
		log.Fatalln("repository.NewAttachment: could not create attachments directory:", err)
	}

	return &Attachment{
		attDir: cfg.DB.Attachments,
		fs:     os.DirFS(cfg.DB.Attachments),
		db:     db,
	}
}

func (a *Attachment) Create(att *domain.Attachment) error {
	err := a.db.Save(att)
	if err != nil {
		return err
	}

	return os.WriteFile(a.getPath(att), att.Data, 0644)
}

// getAttachmentPath returns the path to the attachment file on the file system.
func (a *Attachment) getPath(att *domain.Attachment) string {
	return path.Join(a.attDir, att.File())
}

func (a *Attachment) GetFS() fs.FS {
	return a.fs
}

func (a *Attachment) Get(uuid string) (*domain.Attachment, error) {
	var att domain.Attachment
	err := a.db.One("UUID", uuid, att)
	if err != nil {
		return nil, err
	}

	return &att, nil
}

func (a *Attachment) GetData(att *domain.Attachment) ([]byte, error) {
	data, err := os.ReadFile(a.getPath(att))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (a *Attachment) ListByMessage(msg *domain.Message) ([]domain.Attachment, error) {
	var atts []domain.Attachment
	err := a.db.Select(q.Eq("MessageUUID", msg.UUID)).Find(&atts)
	if err != nil {
		if err == storm.ErrNotFound {
			return []domain.Attachment{}, nil
		}
		return nil, err
	}

	return atts, nil
}

func (a *Attachment) DeleteData(att *domain.Attachment) error {
	return os.Remove(a.getPath(att))
}
