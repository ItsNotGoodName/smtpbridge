package events

import "github.com/ItsNotGoodName/smtpbridge/internal/core"

func OnEnvelopeCreated(app core.App, fn func(cc *core.Context, evt core.EventEnvelopeCreated)) {
	app.Bus.Mutex.Lock()
	app.Bus.EnvelopeCreated = append(app.Bus.EnvelopeCreated, fn)
	app.Bus.Mutex.Unlock()
}

func PublishEnvelopeCreated(cc *core.Context, id int64) {
	cc.Bus.Mutex.Lock()
	for _, v := range cc.Bus.EnvelopeCreated {
		v(cc, core.EventEnvelopeCreated{ID: id})
	}
	cc.Bus.Mutex.Unlock()
}

func OnEnvelopeDeleted(app core.App, fn func(cc *core.Context, evt core.EventEnvelopeDeleted)) {
	app.Bus.Mutex.Lock()
	app.Bus.EnvelopeDeleted = append(app.Bus.EnvelopeDeleted, fn)
	app.Bus.Mutex.Unlock()
}

func PublishEnvelopeDeleted(cc *core.Context, ids ...int64) {
	cc.Bus.Mutex.Lock()
	for _, v := range cc.Bus.EnvelopeDeleted {
		v(cc, core.EventEnvelopeDeleted{IDS: ids})
	}
	cc.Bus.Mutex.Unlock()
}

func OnGardenStart(app core.App, fn func(cc *core.Context, evt core.EventGardenStart)) {
	app.Bus.Mutex.Lock()
	if app.Bus.GardenStart != nil {
		panic("GardenStart handler is being redefined")
	}
	app.Bus.GardenStart = fn
	app.Bus.Mutex.Unlock()
}

func PublishGardenStart(cc *core.Context) <-chan bool {
	res := make(chan bool)
	cc.Bus.Mutex.Lock()
	cc.Bus.GardenStart(cc, core.EventGardenStart{Response: res})
	cc.Bus.Mutex.Unlock()
	return res
}
