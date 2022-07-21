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
	size     int64
	fileName string
	modtime  time.Time
}

type Data struct {
	mu       sync.Mutex
	pool     map[int64]dataBlock
	poolSize int64
	size     int64
	idChan   chan int64
}

func NewData(limit int64, size int64) *Data {
	return &Data{
		idChan: make(chan int64, limit),
		size:   size,
		pool:   make(map[int64]dataBlock),
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
	d.poolSize += size
	// Queue id
	select {
	case d.idChan <- att.ID:
	default:
		id := <-d.idChan
		d.deleteData(id)
		d.idChan <- att.ID
	}

	// Clean up pool if full
	if d.poolSize > d.size {
		for id := range d.idChan {
			d.deleteData(id)
			if d.poolSize <= d.size {
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
	d.poolSize -= data.size

	return nil
}

func (d *Data) DataFS() fs.FS {
	return d
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
