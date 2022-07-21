package memdb

import (
	"context"
	"io/fs"
	"sync"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

const (
	maxDataCount    = 30
	maxDataPoolSize = 1024 * 1024 * 100 // 100 MiB
)

type dataBlock struct {
	data     []byte
	size     int64
	fileName string
	modtime  time.Time
}

type Data struct {
	mu   sync.Mutex
	pool map[int64]dataBlock
	size int64
	idCh chan int64
}

func NewData() *Data {
	return &Data{
		idCh: make(chan int64, maxDataCount),
		pool: make(map[int64]dataBlock),
	}
}

func (d *Data) CreateData(ctx context.Context, att *envelope.Attachment, data []byte) error {
	d.mu.Lock()
	// Don't create data if it already exists
	if _, ok := d.pool[att.ID]; ok {
		d.mu.Unlock()
		return core.ErrDataExists
	}

	// Create data
	size := int64(len(data))
	d.pool[att.ID] = dataBlock{data: data, size: size, fileName: att.FileName(), modtime: time.Now()}
	d.size += size
	// Queue id
	select {
	case d.idCh <- att.ID:
	default:
		id := <-d.idCh
		d.deleteData(id)
		d.idCh <- att.ID
	}

	// Clean up pool if full
	if d.size > maxDataPoolSize {
		for id := range d.idCh {
			d.deleteData(id)
			if d.size <= maxDataPoolSize {
				break
			}
		}
	}
	d.mu.Unlock()

	return nil
}

func (d *Data) GetData(ctx context.Context, att *envelope.Attachment) ([]byte, error) {
	d.mu.Lock()
	block, ok := d.pool[att.ID]
	if !ok {
		d.mu.Unlock()
		return nil, core.ErrDataNotFound
	}
	d.mu.Unlock()

	return block.data, nil
}

func (d *Data) DeleteData(ctx context.Context, att *envelope.Attachment) error {
	d.mu.Lock()
	err := d.deleteData(att.ID)
	d.mu.Unlock()
	return err
}

func (d *Data) deleteData(id int64) error {
	data, ok := d.pool[id]
	if !ok {
		return core.ErrDataNotFound
	}

	delete(d.pool, id)
	d.size -= data.size

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

	d.mu.Lock()
	block, ok := d.pool[id]
	if !ok || block.fileName != name {
		d.mu.Unlock()
		return nil, fs.ErrNotExist
	}
	d.mu.Unlock()

	return newDataFile(block), nil
}
