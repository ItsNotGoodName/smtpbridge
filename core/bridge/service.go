package bridge

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/smtpbridge/core/attachment"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/message"
)

type BridgeService struct {
	bridges            []Bridge
	messageService     message.Service
	endpointRepository endpoint.Repository
	endpointService    endpoint.Service
}

func NewBridgeService(bridges []Bridge, messageSerivce message.Service, endpointRepository endpoint.Repository, endpointService endpoint.Service) *BridgeService {
	return &BridgeService{
		bridges:            bridges,
		messageService:     messageSerivce,
		endpointRepository: endpointRepository,
		endpointService:    endpointService,
	}
}

func (bs *BridgeService) ListByMessage(msg *message.Message) []*Bridge {
	var bridges []*Bridge
	for _, bridge := range bs.bridges {
		if !bridge.Match(msg) {
			continue
		}
		bridges = append(bridges, &bridge)
	}

	return bridges
}

func (bs *BridgeService) HandleMessage(ctx context.Context, bridges []*Bridge, msg *message.Message, atts []attachment.Attachment) error {
	defer bs.messageService.Processed(ctx, msg)

	emsg := endpoint.NewMessage(msg)
	eatts, err := endpoint.NewAttachments(atts)
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
				log.Println("bridge.BridgeService.HandleMessage:", res)
			}
			count++
		}

		log.Printf("bridge.BridgeService.MessageHandle: sent message %d  to %d endpoint(s)", emsg.ID, count)
	}

	return nil
}
