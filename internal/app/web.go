package app

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/web"
)

type WebFileStore struct {
	name string
	url  string
}

func NewWebFileStore(name, url string) WebFileStore {
	return WebFileStore{
		name: name,
		url:  url,
	}
}

func (w WebFileStore) File() (fs.File, error) {
	return web.FS.Open(w.name)
}

func (w WebFileStore) Reader(ctx context.Context, att models.Attachment) (io.ReadCloser, error) {
	return w.File()
}

func (w WebFileStore) Path(ctx context.Context, att models.Attachment) (string, error) {
	if w.url == "" {
		return "", fmt.Errorf("app: url is empty")
	}

	return w.url + "/" + w.name, nil
}
