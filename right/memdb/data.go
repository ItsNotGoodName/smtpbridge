package memdb

import (
	"context"
	"io/fs"
	"sync"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type dataBlock struct {
	data     []byte
	length   int64
	fileName string
	modtime  time.Time
}

type Data struct {
	blockMu sync.Mutex
	block   map[int64]dataBlock
}

func NewData() *Data {
	return &Data{
		block: make(map[int64]dataBlock),
	}
}

func (d *Data) CreateData(ctx context.Context, att *envelope.Attachment, data []byte) error {
	d.blockMu.Lock()
	d.block[att.ID] = dataBlock{data, int64(len(data)), att.FileName(), time.Now()}
	delete(d.block, att.ID-maxData)
	d.blockMu.Unlock()

	return nil
}

func (d *Data) GetData(ctx context.Context, att *envelope.Attachment) ([]byte, error) {
	d.blockMu.Lock()
	block, ok := d.block[att.ID]
	if !ok {
		d.blockMu.Unlock()
		return nil, core.ErrDataNotFound
	}
	d.blockMu.Unlock()

	return block.data, nil
}

func (d *Data) DeleteData(ctx context.Context, att *envelope.Attachment) error {
	d.blockMu.Lock()
	delete(d.block, att.ID)
	d.blockMu.Unlock()

	return nil
}

func (d *Data) DataFS() (fs.FS, error) {
	return d, nil
}

func (d *Data) Open(name string) (fs.File, error) {
	id, err := envelope.AttachmentIDFromFileName(name)
	if err != nil {
		return nil, fs.ErrNotExist
	}

	d.blockMu.Lock()
	block, ok := d.block[id]
	if !ok || block.fileName != name {
		d.blockMu.Unlock()
		return nil, fs.ErrNotExist
	}
	d.blockMu.Unlock()

	return newDataFile(block), nil
}
