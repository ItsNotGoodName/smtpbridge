package bridge

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type BridgeService struct {
	bridges         []Bridge
	messageService  message.Service
	endpointService endpoint.Service
}

func NewBridgeService(bridges []Bridge, messageSerivce message.Service, endpointService endpoint.Service) *BridgeService {
	return &BridgeService{
		bridges:         bridges,
		messageService:  messageSerivce,
		endpointService: endpointService,
	}
}

func (bs *BridgeService) ListByEnvelope(env envelope.Envelope) []*Bridge {
	var bridges []*Bridge
	for _, bridge := range bs.bridges {
		if !bridge.Match(env) {
			continue
		}
		bridges = append(bridges, &bridge)
	}

	return bridges
}

func (bs *BridgeService) HandleEnvelope(ctx context.Context, bridges []*Bridge, env envelope.Envelope) error {
	defer bs.messageService.Processed(ctx, env.Message)

	emsg := endpoint.NewMessage(env.Message)
	eatts, err := endpoint.NewAttachments(env.Attachments)
	if err != nil {
		return err
	}

	for _, bd := range bridges {
		var req []endpoint.SendRequest
		for _, f := range bd.Endpoints {
			req = append(req, endpoint.SendRequest{
				Envelope: f.Envelope(emsg, eatts),
				Facade:   f.Facade,
			})
		}

		var count int
		resC := bs.endpointService.Send(ctx, req)
		for res := range resC {
			if res.Error != nil {
				log.Println("bridge.BridgeService.HandleEnvelope:", res)
			}
			count++
		}

		log.Printf("bridge.BridgeService.MessageEnvelope: sent envelope %d to %d endpoint(s)", emsg.ID, count)
	}

	return nil
}
