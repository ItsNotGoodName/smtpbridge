package bridge

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/event"
)

type BridgeService struct {
	eventPub        *event.Pub
	envelopeService envelope.Service
	endpointService endpoint.Service
	bridgesMu       sync.Mutex
	bridges         []Bridge
}

func NewBridgeService(pub *event.Pub, envelopeService envelope.Service, endpointService endpoint.Service) *BridgeService {
	return &BridgeService{
		eventPub:        pub,
		envelopeService: envelopeService,
		endpointService: endpointService,
	}
}

func (bs *BridgeService) ListBridge() []Bridge {
	bridges := []Bridge{}

	bs.bridgesMu.Lock()
	bridges = append(bridges, bs.bridges...)
	bs.bridgesMu.Unlock()

	return bridges
}

func (bs *BridgeService) CreateBridge(req *CreateBridgeRequest) error {
	for _, end := range req.Endpoints {
		if _, err := bs.endpointService.GetEndpoint(end); err != nil {
			return fmt.Errorf("%w: '%s'", err, end)
		}
	}

	filter, err := NewFilter(req.To, req.From, req.ToRegex, req.FromRegex, req.MatchTemplate)
	if err != nil {
		return err
	}

	bridge := NewBridge(filter, req.Endpoints)
	bs.bridgesMu.Lock()
	bs.bridges = append(bs.bridges, bridge)
	bs.bridgesMu.Unlock()

	return nil
}

func (bs *BridgeService) send(ctx context.Context, env *envelope.Envelope) error {
	// Convert envelope attachments to endpoint attachments
	atts := []endpoint.Attachment{}
	for _, att := range env.Attachments {
		data, err := bs.envelopeService.GetData(ctx, &att)
		if err != nil {
			if errors.Is(err, core.ErrDataNotFound) {
				log.Println("bridge.BridgeService.send:", err)
				continue
			}
			return err
		}

		atts = append(atts, endpoint.NewAttachment(&att, data))
	}

	for _, brid := range bs.ListBridge() {
		// Match bridge
		if !brid.Filter.Match(env) {
			continue
		}

		// Send to all endpoints
		if len(brid.Endpoints) == 0 {
			for _, end := range bs.endpointService.ListEndpoint() {
				text, err := end.Text(env)
				if err != nil {
					log.Println("bridge.BridgeService.send:", err)
					continue
				}

				end.Sender.Send(ctx, text, atts)
			}

			return nil
		}

		// Send to endpoints
		for _, endpoitName := range brid.Endpoints {
			end, err := bs.endpointService.GetEndpoint(endpoitName)
			if err != nil {
				log.Println("bridge.BridgeService.Run:", err)
				continue
			}

			text, err := end.Text(env)
			if err != nil {
				log.Println("bridge.BridgeService.Run:", err)
				continue
			}

			end.Sender.Send(ctx, text, atts)
		}
	}

	return nil
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
			if err := bs.send(ctx, event.Data.(*envelope.Envelope)); err != nil {
				log.Println("bridge.BridgeService.Run:", err)
			}
		}
	}
}
