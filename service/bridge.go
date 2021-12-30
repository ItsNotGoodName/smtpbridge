package service

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/domain"
)

type Bridge struct {
	bridges []domain.Bridge
}

func NewBridge(dao domain.DAO, bridges []domain.Bridge) *Bridge {
	for _, bridge := range bridges {
		for _, endpoint := range bridge.Endpoints {
			if _, err := dao.Endpoint.Get(endpoint); err != nil {
				log.Fatalln("service.NewBridge:", err)
			}
		}
	}

	return &Bridge{bridges: bridges}
}

func (b *Bridge) GetBridges(msg *domain.Message) []domain.Bridge {
	var bridges []domain.Bridge
	for _, bridge := range b.bridges {
		if !bridge.Match(msg) {
			continue
		}
		bridges = append(bridges, bridge)
	}

	return bridges
}
