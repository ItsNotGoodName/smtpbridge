package app

import "io/fs"

func (a *App) AttachmentFS() fs.FS {
	return a.dataRepository.FS()
}
