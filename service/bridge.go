package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Bridge struct {
	// TODO: pointer to bridge
	bridges []domain.Bridge
}

func NewBridge(cfg *config.Config, endpointREPO domain.EndpointRepositoryPort) *Bridge {
	// TODO: move somewhere else
	for _, bridge := range cfg.Bridges {
		for _, endpoint := range bridge.Endpoints {
			if _, err := endpointREPO.Get(endpoint); err != nil {
				log.Fatalln("service.NewBridge:", err)
			}
		}
	}

	return &Bridge{bridges: cfg.Bridges}
}

func (b *Bridge) ListByMessage(msg *domain.Message) []domain.Bridge {
	var bridges []domain.Bridge
	for _, bridge := range b.bridges {
		if !bridge.Match(msg) {
			continue
		}
		bridges = append(bridges, bridge)
	}

	return bridges
}
