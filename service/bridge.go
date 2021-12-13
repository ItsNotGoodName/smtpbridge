package service

import (
	"github.com/ItsNotGoodName/go-smtpbridge/app"
)

type Bridge struct {
	endpoints map[string]app.EndpointPort
	bridges   []app.Bridge
}

func NewBridge(bridges []app.Bridge, endpoints map[string]app.EndpointPort) *Bridge {
	return &Bridge{endpoints: endpoints, bridges: bridges}
}

func (b *Bridge) GetEndpoints(msg *app.Message) []app.EndpointPort {
	var endpoints []app.EndpointPort
	for _, bridge := range b.bridges {
		if !bridge.Match(msg) {
			continue
		}

		for _, endpointName := range bridge.Endpoints {
			endpoint, ok := b.endpoints[endpointName]
			if !ok {
				panic("endpoint not found")
			}
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}
