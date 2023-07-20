package files

import (
	"io/ioutil"
	"os"
	"path"
	"sort"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
)

func CreateFile(cc *core.Context, datt envelope.DataAttachment) error {
	return os.WriteFile(filePath(cc, datt.Attachment), datt.Data, 0644)
}

func GetFileData(cc *core.Context, att *envelope.Attachment) (envelope.DataAttachment, error) {
	data, err := os.ReadFile(filePath(cc, att))
	return envelope.DataAttachment{Attachment: att, Data: data}, err
}

func GetFile(cc *core.Context, att *envelope.Attachment) (*os.File, error) {
	return os.Open(filePath(cc, att))
}

func DeleteFile(cc *core.Context, att *envelope.Attachment) error {
	return os.Remove(filePath(cc, att))
}

func Size(cc *core.Context) (int64, error) {
	files, err := ioutil.ReadDir(cc.File.Dir)
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

func DeleteFileUntilSize(cc *core.Context, currentSize, maxSize int64) error {
	files, err := ioutil.ReadDir(cc.File.Dir)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	for i := range files {
		if currentSize < maxSize {
			break
		}

		if err := os.Remove(path.Join(cc.File.Dir, files[i].Name())); err != nil {
			return err
		}

		currentSize -= files[i].Size()
	}

	return nil
}

func filePath(cc *core.Context, att *envelope.Attachment) string {
	return path.Join(cc.File.Dir, att.FileName())
}
