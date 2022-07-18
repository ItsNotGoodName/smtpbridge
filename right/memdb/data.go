package memdb

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type Data struct {
	mu   sync.Mutex
	data map[int64][]byte
}

func NewData() *Data {
	return &Data{
		data: make(map[int64][]byte),
	}
}

func (d *Data) CreateData(ctx context.Context, att *envelope.Attachment, data []byte) error {
	d.mu.Lock()
	d.data[att.ID] = data
	// Delete oldest attachment data if full
	if len(data) > maxData {
		delete(d.data, att.ID-maxData)
	}
	d.mu.Unlock()

	return nil
}

func (d *Data) GetData(ctx context.Context, att *envelope.Attachment) ([]byte, error) {
	d.mu.Lock()
	data, ok := d.data[att.ID]
	if !ok {
		d.mu.Unlock()
		return nil, core.ErrDataNotFound
	}
	d.mu.Unlock()

	return data, nil
}

func (d *Data) DeleteData(ctx context.Context, att *envelope.Attachment) error {
	d.mu.Lock()
	delete(d.data, att.ID)
	d.mu.Unlock()

	return nil
}
