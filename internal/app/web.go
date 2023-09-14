package app

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/web"
)

// WebTestFileStore is used for testing endpoints.
type WebTestFileStore struct {
	name string
	url  string
}

func NewWebTestFileStore(name, url string) WebTestFileStore {
	return WebTestFileStore{
		name: name,
		url:  url,
	}
}

func (w WebTestFileStore) File() (fs.File, error) {
	return web.FS.Open(w.name)
}

func (w WebTestFileStore) Reader(ctx context.Context, att models.Attachment) (io.ReadCloser, error) {
	return w.File()
}

func (w WebTestFileStore) Path(ctx context.Context, att models.Attachment) (string, error) {
	if w.url == "" {
		return "", fmt.Errorf("app: url is empty")
	}

	return w.url + "/" + w.name, nil
}
