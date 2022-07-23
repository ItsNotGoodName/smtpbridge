package bridge

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
)

type (
	Service interface{}
)

type BridgeService struct {
	eventPub        *event.Pub
	envelopeService envelope.Service
	endpointService endpoint.Service
}

func NewBridgeService(pub *event.Pub, envelopeService envelope.Service, endpointService endpoint.Service) *BridgeService {
	return &BridgeService{
		eventPub:        pub,
		envelopeService: envelopeService,
		endpointService: endpointService,
	}
}

func (bs *BridgeService) Run(ctx context.Context, doneC chan<- struct{}) {
	log.Println("bridge.BridgeService.Run: started")

	eventChan := make(chan event.Event, 100)

	bs.eventPub.Subscribe(event.TopicEnvelopeCreated, eventChan)
	for {
		select {
		case <-ctx.Done():
			log.Println("bridge.BridgeService.Run: stopped")
			doneC <- struct{}{}
			return
		case event := <-eventChan:
			env := event.Data.(*envelope.Envelope)

			// Convert envelope attachments to endpoint attachments
			atts := []endpoint.Attachment{}
			for _, att := range env.Attachments {
				data, err := bs.envelopeService.GetData(ctx, &att)
				if err != nil && err != core.ErrDataNotFound {
					log.Println("bridge.BridgeService.Run:", err)
					continue
				}

				atts = append(atts, endpoint.NewAttachment(&att, data))
			}

			ends := bs.endpointService.ListEndpoint()

			// Send to all endpoints
			for _, end := range ends {
				text, err := end.Text(env)
				if err != nil {
					log.Println("bridge.BridgeService.Run:", err)
				} else {
					end.Sender.Send(ctx, text, atts)
				}
			}
		}
	}
}
