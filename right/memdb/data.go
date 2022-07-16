package memdb

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
)

type Data struct {
	dataMu sync.Mutex
	data   map[int64][]byte
}

func NewData() *Data {
	return &Data{
		data: make(map[int64][]byte),
	}
}

func (d *Data) CreateData(ctx context.Context, att *envelope.Attachment, data []byte) error {
	d.dataMu.Lock()
	d.data[att.ID] = data
	d.dataMu.Unlock()

	return nil
}

func (d *Data) GetData(ctx context.Context, att *envelope.Attachment) ([]byte, error) {
	d.dataMu.Lock()
	data, ok := d.data[att.ID]
	if !ok {
		d.dataMu.Unlock()
		return nil, core.ErrDataNotFound
	}
	d.dataMu.Unlock()

	return data, nil
}

func (d *Data) DeleteData(ctx context.Context, att *envelope.Attachment) error {
	d.dataMu.Lock()
	delete(d.data, att.ID)
	d.dataMu.Unlock()

	return nil
}
