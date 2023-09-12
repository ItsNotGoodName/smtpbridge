package app

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/web"
)

type TestFileStore struct {
	name string
	url  string
}

func NewTestFileStore(name, url string) TestFileStore {
	return TestFileStore{
		name: name,
		url:  url,
	}
}

func (t TestFileStore) File() (fs.File, error) {
	return web.FS.Open(t.name)
}

func (t TestFileStore) Reader(ctx context.Context, att models.Attachment) (io.ReadCloser, error) {
	return t.File()
}

func (t TestFileStore) Path(ctx context.Context, att models.Attachment) (string, error) {
	if t.url == "" {
		return "", fmt.Errorf("app: url is empty")
	}

	return t.url + "/" + t.name, nil
}
