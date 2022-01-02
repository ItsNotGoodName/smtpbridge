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

type attachmentModel struct {
	UUID        string                `storm:"id"`
	Name        string                ``
	Type        domain.AttachmentType ``
	MessageUUID string                ``
}

func convertAttachmentD(att *domain.Attachment) *attachmentModel {
	return &attachmentModel{
		UUID:        att.UUID,
		Name:        att.Name,
		Type:        att.Type,
		MessageUUID: att.MessageUUID,
	}
}

func convertAttachmentM(attM *attachmentModel) *domain.Attachment {
	return &domain.Attachment{
		UUID:        attM.UUID,
		Name:        attM.Name,
		Type:        attM.Type,
		MessageUUID: attM.MessageUUID,
	}
}

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
	err := a.db.Save(convertAttachmentD(att))
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
	var attM *attachmentModel
	err := a.db.One("UUID", uuid, attM)
	if err != nil {
		return nil, err
	}

	return convertAttachmentM(attM), nil
}

func (a *Attachment) GetData(att *domain.Attachment) ([]byte, error) {
	data, err := os.ReadFile(a.getPath(att))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (a *Attachment) ListByMessage(msg *domain.Message) ([]domain.Attachment, error) {
	var attsM []attachmentModel
	err := a.db.Select(q.Eq("MessageUUID", msg.UUID)).Find(&attsM)
	if err != nil {
		if err == storm.ErrNotFound {
			return []domain.Attachment{}, nil
		}
		return nil, err
	}

	var atts []domain.Attachment
	for _, attM := range attsM {
		atts = append(atts, *convertAttachmentM(&attM))
	}

	return atts, nil
}

func (a *Attachment) DeleteData(att *domain.Attachment) error {
	return os.Remove(a.getPath(att))
}
