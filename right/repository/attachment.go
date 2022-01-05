package repository

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

type attachmentModel struct {
	UUID        string              `json:"uuid" storm:"id"`
	Name        string              `json:"name"`
	Type        core.AttachmentType `json:"type"`
	MessageUUID string              `json:"message_uuid"`
}

func convertAttachmentD(att *core.Attachment) *attachmentModel {
	return &attachmentModel{
		UUID:        att.UUID,
		Name:        att.Name,
		Type:        att.Type,
		MessageUUID: att.MessageUUID,
	}
}

func convertAttachmentM(attM *attachmentModel) *core.Attachment {
	return &core.Attachment{
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

func NewAttachment(cfg *config.Config, db *Database) *Attachment {
	err := os.MkdirAll(cfg.DB.Attachments, 0755)
	if err != nil {
		log.Fatalln("repository.NewAttachment: could not create attachments directory:", err)
	}

	return &Attachment{
		attDir: cfg.DB.Attachments,
		fs:     os.DirFS(cfg.DB.Attachments),
		db:     db.db,
	}
}

func (a *Attachment) Create(att *core.Attachment) error {
	err := a.db.Save(convertAttachmentD(att))
	if err != nil {
		return err
	}

	return os.WriteFile(a.getPath(att), att.Data, 0644)
}

func (a *Attachment) Count() (int, error) {
	count, err := a.db.Count(&attachmentModel{})
	if err == storm.ErrNotFound {
		return 0, nil
	}

	return count, nil
}

func (a *Attachment) CountByMessage(msg *core.Message) (int, error) {
	// TODO: do i use MessageUUID or message_uuid?
	count, err := a.db.Select(q.Eq("MessageUUID", msg.UUID)).Count(&attachmentModel{})
	if err == storm.ErrNotFound {
		return 0, nil
	}

	return count, nil
}

// getAttachmentPath returns the path to the attachment file on the file system.
func (a *Attachment) getPath(att *core.Attachment) string {
	return path.Join(a.attDir, att.File())
}

func (a *Attachment) GetFS() fs.FS {
	return a.fs
}

func (a *Attachment) Get(uuid string) (*core.Attachment, error) {
	var attM *attachmentModel
	err := a.db.One("UUID", uuid, attM)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, core.ErrAttachmentNotFound
		}
		return nil, err
	}

	return convertAttachmentM(attM), nil
}

func (a *Attachment) GetSizeAll() (int64, error) {
	if err := os.Chdir(a.attDir); err != nil {
		return 0, err
	}

	files, err := ioutil.ReadDir(a.attDir)
	if err != nil {
		return 0, err
	}

	dirSize := int64(0)
	for _, file := range files {
		if file.Mode().IsRegular() {
			dirSize += file.Size()
		}
	}

	return dirSize, nil
}

func (a *Attachment) ListByMessage(msg *core.Message) ([]core.Attachment, error) {
	var attsM []attachmentModel
	// TODO: do i use MessageUUID or message_uuid?
	err := a.db.Select(q.Eq("MessageUUID", msg.UUID)).Find(&attsM)
	if err != nil {
		if err == storm.ErrNotFound {
			return []core.Attachment{}, nil
		}
		return nil, err
	}

	var atts []core.Attachment
	for _, attM := range attsM {
		atts = append(atts, *convertAttachmentM(&attM))
	}

	return atts, nil
}

func (a *Attachment) LoadData(att *core.Attachment) error {
	data, err := os.ReadFile(a.getPath(att))
	if err != nil {
		if err == os.ErrNotExist {
			return core.ErrAttachmentNotFound
		}
		return err
	}

	att.Data = data

	return nil
}

func (a *Attachment) deleteData(att *core.Attachment) error {
	return os.Remove(a.getPath(att))
}
