package filedb

import (
	"context"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
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

func (d *Data) Create(ctx context.Context, att *attachment.Attachment) error {
	var (
		data []byte
		err  error
	)
	if data, err = att.GetData(); err != nil {
		return err
	}

	return os.WriteFile(d.getPath(att), data, 0644)
}

func (d *Data) FS() fs.FS {
	return d.fs
}

func (d *Data) Delete(ctx context.Context, att *attachment.Attachment) error {
	err := os.Remove(d.getPath(att))
	if err == os.ErrNotExist {
		return attachment.ErrNotFound
	}
	return err
}

func (d *Data) Load(ctx context.Context, att *attachment.Attachment) error {
	data, err := os.ReadFile(d.getPath(att))
	if err != nil {
		if err == os.ErrNotExist {
			return attachment.ErrNotFound
		}
		return err
	}

	if err := att.SetData(data); err != nil {
		return err
	}

	return nil
}

func (d *Data) Size(ctx context.Context) (int64, error) {
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

// getPath returns the path to the attachment file on the file system.
func (d *Data) getPath(att *attachment.Attachment) string {
	return path.Join(d.dir, att.File())
}
