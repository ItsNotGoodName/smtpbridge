package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

type Bridge struct {
	bridges []app.Bridge
}

func NewBridge(endpointREPO app.EndpointRepositoryPort, bridges []app.Bridge) *Bridge {
	for _, bridge := range bridges {
		for _, endpoint := range bridge.Endpoints {
			if _, err := endpointREPO.Get(endpoint); err != nil {
				log.Fatalln("service.NewBridge:", err)
			}
		}
	}

	return &Bridge{bridges: bridges}
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
