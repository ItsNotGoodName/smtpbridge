package filedb

import (
	"context"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type Data struct {
	dir string
	fs  fs.FS
}

func NewData(dir string) *Data {
	return &Data{
		dir: dir,
		fs:  os.DirFS(dir),
	}
}

func (d *Data) CreateData(ctx context.Context, att *envelope.Attachment, data []byte) error {
	return os.WriteFile(d.path(att), data, 0644)
}

func (d *Data) GetData(ctx context.Context, att *envelope.Attachment) ([]byte, error) {
	data, err := os.ReadFile(d.path(att))
	return data, checkDataError(err)
}

func (d *Data) DeleteData(ctx context.Context, att *envelope.Attachment) error {
	return checkDataError(os.Remove(d.path(att)))
}

func (d *Data) DataSize(ctx context.Context) (int64, error) {
	if err := os.Chdir(d.dir); err != nil {
		return 0, err
	}

	files, err := ioutil.ReadDir(d.dir)
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

func (d *Data) DataFS() fs.FS {
	return d.fs
}

func (d *Data) path(att *envelope.Attachment) string {
	return path.Join(d.dir, att.FileName())
}

func checkDataError(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return core.ErrDataNotFound
	}

	return err
}
