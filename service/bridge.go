package service

import (
	"fmt"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Bridge struct {
	endpoints map[string]app.EndpointPort
	bridges   []app.Bridge
}

func NewBridge(bridges []app.Bridge, endpoints map[string]app.EndpointPort) *Bridge {
	return &Bridge{endpoints: endpoints, bridges: bridges}
}

func (b *Bridge) GetBridges(msg *app.Message) []app.Bridge {
	var bridges []app.Bridge
	for _, bridge := range b.bridges {
		if !bridge.Match(msg) {
			continue
		}
		bridges = append(bridges, bridge)
	}

	return bridges
}

func (b *Bridge) GetEndpoint(name string) app.EndpointPort {
	endpoint, ok := b.endpoints[name]
	if !ok {
		panic(fmt.Sprintf("endpoint '%s' not found", name))
	}

	return endpoint
}
