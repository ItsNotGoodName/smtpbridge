package app

import "io/fs"

type Attachment struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	File string `json:"path"`
}

func (a *App) AttachmentGetFS() fs.FS {
	return a.attachmentREPO.GetFS()
}
