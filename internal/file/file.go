package file

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

type Store struct {
	directory string
}

// Create implements core.FileStore.
func (s Store) Create(ctx context.Context, att models.Attachment, data io.Reader) error {
	f, err := os.Create(s.filePath(att))
	if err != nil {
		return err
	}
	_, err = io.Copy(f, data)
	return err
}

func NewStore(directory string) Store {
	return Store{
		directory: directory,
	}
}

// Remove implements core.FileStore.
func (s Store) Remove(ctx context.Context, att models.Attachment) error {
	err := os.Remove(s.filePath(att))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

// Path implements core.FileStore.
func (s Store) Path(ctx context.Context, att models.Attachment) (string, error) {
	return s.filePath(att), nil
}

// Reader implements core.FileStore.
func (s Store) Reader(ctx context.Context, att models.Attachment) (io.ReadCloser, error) {
	return os.Open(s.filePath(att))
}

// Size implements core.FileStore.
func (s Store) Size(ctx context.Context) (int64, error) {
	files, err := ioutil.ReadDir(s.directory)
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

// Trim implements core.FileStore.
func (s Store) Trim(ctx context.Context, size int64, minAge time.Time) (int, error) {
	currentSize, err := s.Size(ctx)
	if err != nil {
		return 0, err
	}

	files, err := ioutil.ReadDir(s.directory)
	if err != nil {
		return 0, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	var count int
	for i := range files {
		if currentSize < size {
			break
		}
		if files[i].ModTime().After(minAge) {
			continue
		}

		if err := os.Remove(path.Join(s.directory, files[i].Name())); err != nil {
			return 0, err
		}

		currentSize -= files[i].Size()
		count += 1
	}

	return count, nil
}

func (s Store) filePath(att models.Attachment) string {
	return path.Join(s.directory, att.FileName())
}

func (s Store) Open(name string) (fs.File, error) {
	return os.Open(path.Join(s.directory, name))
}
