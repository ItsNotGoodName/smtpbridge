package events

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

func OnStorageRead(app core.App, fn func(cc *core.Context, evt core.EventStorageRead)) {
	app.Bus.Mutex.Lock()
	app.Bus.StorageRead = append(app.Bus.StorageRead, fn)
	app.Bus.Mutex.Unlock()
}

func PublishStorageRead(cc *core.Context, storage models.Storage) {
	cc.Bus.Mutex.Lock()
	for _, v := range cc.Bus.StorageRead {
		v(cc, core.EventStorageRead{Storage: storage})
	}
	cc.Bus.Mutex.Unlock()
}
