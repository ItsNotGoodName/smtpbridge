package app

import "io/fs"

func (a *App) AttachmentGetFS() fs.FS {
	return a.dao.Attachment.GetAttachmentFS()
}
