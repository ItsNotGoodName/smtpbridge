package filedb

import (
	"context"
	"fmt"
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

func (d *Data) Create(ctx context.Context, att *attachment.Attachment, data []byte) error {
	return os.WriteFile(d.path(att), data, 0644)
}

func (d *Data) FS() fs.FS {
	return d.fs
}

func (d *Data) Delete(ctx context.Context, att *attachment.Attachment) error {
	err := os.Remove(d.path(att))
	if err == os.ErrNotExist {
		return attachment.ErrNotFound
	}

	return err
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

func (d *Data) Get(ctx context.Context, att *attachment.Attachment) ([]byte, error) {
	return ioutil.ReadFile(d.path(att))
}

func (d *Data) Remote() bool {
	return false
}

func (d *Data) URL(att *attachment.Attachment) string {
	return ""
}

func (d *Data) File(att *attachment.Attachment) string {
	return fmt.Sprintf("%d.%s", att.ID, att.Type)
}

// path returns the path to the attachment file on the file system.
func (d *Data) path(att *attachment.Attachment) string {
	return path.Join(d.dir, d.File(att))
}
