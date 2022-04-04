package app

import "io/fs"

func (a *App) AttachmentDataFS() fs.FS {
	return a.attachmentDataService.FS()
}

func (a *App) AttachmentDataURI() string {
	return a.attachmentDataService.URI()
}

func (a *App) AttachmentDataRemote() bool {
	return a.attachmentDataService.Remote()
}
